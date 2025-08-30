package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
//要求 ：
//编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

type Account struct {
	ID      uint `gorm:"primaryKey"`
	Balance float64
}

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

func main() {
	dsn := "root:root@tcp(10.4.8.22:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移表结构
	db.AutoMigrate(&Account{}, &Transaction{})

	// 示例账户ID
	fromAccountID := uint(1)
	toAccountID := uint(2)
	amount := 100.0

	// 执行转账
	err = TransferMoney(db, fromAccountID, toAccountID, amount)
	if err != nil {
		log.Println("转账失败:", err)
	} else {
		fmt.Println("转账成功")
	}
}

func TransferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 检查转出账户余额
		var fromAccount Account
		if err := tx.First(&fromAccount, fromID).Error; err != nil {
			return fmt.Errorf("查询转出账户失败: %v", err)
		}

		if fromAccount.Balance < amount {
			return fmt.Errorf("账户余额不足")
		}

		// 检查转入账户是否存在
		var toAccount Account
		if err := tx.First(&toAccount, toID).Error; err != nil {
			return fmt.Errorf("查询转入账户失败: %v", err)
		}

		// 更新账户余额
		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-amount).Error; err != nil {
			return fmt.Errorf("扣除转出账户金额失败: %v", err)
		}

		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+amount).Error; err != nil {
			return fmt.Errorf("增加转入账户金额失败: %v", err)
		}

		// 记录交易
		transaction := Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("记录交易失败: %v", err)
		}

		return nil
	})
}
