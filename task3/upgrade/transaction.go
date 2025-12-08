/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-05 13:46:00
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-05 15:46:18
 */
package upgrade

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

/**
题目2：事务语句
假设有两个表：
 accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
 （包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

// 账户表
type account struct {
	Id      uint `gorm: "primaryKey", "autoIncrement"`
	Balance float64
}

// 转账表
type transaction struct {
	Id            uint `gorm: "primaryKey", "autoIncrement"`
	FromAccountId uint
	ToAccountId   uint
	Amount        float64
}

func Charge(db *gorm.DB, accountId uint, amount float64) error {

	errs := db.AutoMigrate(&account{})
	if errs != nil {
		panic("自动迁移失败: " + errs.Error())
	}
	errs1 := db.AutoMigrate(&transaction{})
	if errs1 != nil {
		panic("自动迁移失败: " + errs1.Error())
	}

	var acc account
	err := db.Where("id = ?", accountId).First(&acc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//账户不存在，创建一个新账户
		acc = account{
			Id: accountId,
		}
	}
	acc.Balance += amount
	err1 := db.Save(&acc).Error
	if err1 != nil {
		return err1
	}
	fmt.Println("充值成功:账户：", accountId, "余额：", acc.Balance)
	return nil
}

// 注：事务中所有的ddl操作都要捕获返回信息，不然会静默
func Transfer(db *gorm.DB, fromAccountId uint, toAccountId uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		//转账之前先查询是否有足够的余额
		var acc account
		err := tx.Where("id = ?", fromAccountId).First(&acc).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("A账户不存在")
		}
		if acc.Balance < amount {
			return errors.New("余额不足")
		}
		//校验B账户是否存在
		var bcc account
		err1 := tx.Where("id = ?", toAccountId).First(&bcc).Error
		if errors.Is(err1, gorm.ErrRecordNotFound) {
			return errors.New("B账户不存在")
		}

		//开始转账
		acc.Balance -= amount
		bcc.Balance += amount
		//记录转账信息
		trans := transaction{
			FromAccountId: fromAccountId,
			ToAccountId:   toAccountId,
			Amount:        amount,
		}
		err2 := tx.Create(&trans).Error
		if err2 != nil {
			return err2
		}

		//更新账户transfer
		tx.Save(&acc)
		tx.Save(&bcc)

		fmt.Println("转账成功")
		return nil
	})

}
