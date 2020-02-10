package main

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

type Goods struct {
	Id    int `xorm:"not null pk autoincr INTEGER"`
	Price uint
}

func (self *Goods) BeforeInsert() {
	fmt.Println("before insert good %v", self.Id)
}

func (self *Goods) AfterInsert() {
	fmt.Println("after insert good %v", self.Id)
}

func (self *Goods) BeforeUpdate() {
	fmt.Println("before update good %v", self.Id)
}

func (self *Goods) AfterUpdate() {
	fmt.Println("after update good %v", self.Id)
}

func (self *Goods) BeforeDelete() {
	fmt.Println("after delete good %v", self.Id)
}
func (self *Goods) AfterDelete() {
	fmt.Println("after delete good %v", self.Id)
}

func (self *Goods) BeforeSet(name string, cell xorm.Cell) {
	fmt.Println("before set %v as %v", name, *cell)
}
func (self *Goods) AfterSet(name string, cell xorm.Cell) {
	fmt.Println("after set %v as %v", name, *cell)
}

func sync_table() {
	db, err := xorm.NewEngine("postgres", "postgres://postgres:postgres@localhost:5432/test?sslmode=disable")
	if err != nil {
		fmt.Printf("%v \n", err)
	} else {
		defer db.Close()
		db.ShowSQL(true)
		var ok bool
		ok, err = db.IsTableExist("goods")
		if err != nil {
			fmt.Println("table goods IsTableExist error", err)
		} else {
			fmt.Println("table goods is exist :", ok)
			err = db.Sync2(&Goods{})
			if err != nil {
				fmt.Println("table goods Sync2 error", err)
			} else {
				ok, err = db.IsTableExist("goods")
				if err != nil {
					fmt.Println("table goods IsTableExist error", err)
				} else {
					fmt.Println("table goods is exist :", ok)
					err = db.DropTables("goods")
					if err != nil {
						fmt.Println("table goods drop error:", err)
					} else {
						fmt.Println("table drop ok")
						ok, err = db.IsTableExist("goods")
						if err != nil {
							fmt.Println("table goods IsTableExist error", err)
						} else {
							fmt.Println("table goods is exist :", ok)
						}
					}
				}
			}

		}
	}
}
func insert_table(db *xorm.Engine) {
	good := &Goods{
		Id:    4,
		Price: 0,
	}
	affected, err := db.Insert(good)
	if err != nil {
		fmt.Println("insert err:", err)
	} else {
		fmt.Println("insert affected:", affected)
	}
	goods := []*Goods{{
		Id:    1,
		Price: 3,
	}, {
		Id:    2,
		Price: 2,
	}, {
		Id:    3,
		Price: 1,
	},
	}
	affected, err = db.Insert(goods)
	if err != nil {
		fmt.Println("insert err:", err)
	} else {
		fmt.Println("insert affected:", affected)
	}
}

func update_table(db *xorm.Engine) {
	fmt.Println("--------------------update by cols----------------------")
	good := new(Goods)
	good.Price = 15
	affected, err := db.Id(1).Cols("price").Update(good)
	if err != nil {
		fmt.Println("update err:", err)
	} else {
		fmt.Println("update affected:", affected)
	}
	fmt.Println("--------------------update by map----------------------")
	affected, err = db.Table(new(Goods)).Id(2).Update(map[string]interface{}{"price": 20})
	if err != nil {
		fmt.Println("update err:", err)
	} else {
		fmt.Println("update affected:", affected)
	}

}

func select_get_table(db *xorm.Engine) {
	fmt.Println("--------------------select get----------------------")
	good := new(Goods)
	has, err := db.After(func(bean interface{}) {
		temp := bean.(*Goods)
		fmt.Println("after get select get table instance: ", temp)
	}).Id(3).Get(good)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		if has {
			fmt.Println("good price:", good.Price)
		} else {
			fmt.Println("good not exist")
		}
	}
}
func select_exist_table(db *xorm.Engine) {
	fmt.Println("--------------------select exist----------------------")
	good := new(Goods)
	has, err := db.Id(3).Exist(good)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		if has {
			fmt.Println("good exist")
		} else {
			fmt.Println("good not exist")
		}
	}
}

func select_find_table(db *xorm.Engine) {
	fmt.Println("--------------------select find----------------------")
	goods := make([]Goods, 0)
	err := db.Find(&goods)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, v := range goods {
			fmt.Println("good Price:", v.Price)
		}
	}
}

func select_row_table(db *xorm.Engine) {
	fmt.Println("--------------------select row----------------------")
	good := new(Goods)
	rows, err := db.Where("id >?", 1).Rows(good)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(good)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("good price:", good.Price)
		}
	}
}
func select_iterate_table(db *xorm.Engine) {
	fmt.Println("--------------------select iterate----------------------")
	err := db.Where("id >?", 1).Iterate(new(Goods), func(i int, bean interface{}) error {
		good := bean.(*Goods)
		fmt.Println("good price:", good.Price)
		return nil
	})
	if err != nil {
		fmt.Println("error:", err)
	}
}
func sum_table(db *xorm.Engine) {
	fmt.Println("--------------------sum ----------------------")
	total, err := db.Where("id >?", 1).Sum(new(Goods), "price")
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("total price :", total)
	}
	fmt.Println("--------------------sumint ----------------------")
	totalint, err1 := db.Where("id >?", 1).SumInt(new(Goods), "price")
	if err != nil {
		fmt.Println("error:", err1)
	} else {
		fmt.Println("total int price :", totalint)
	}
}
func count_table(db *xorm.Engine) {
	fmt.Println("--------------------count ----------------------")
	count, err := db.Where("id >?", 1).Count(new(Goods))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("count:", count)
	}
}

func query_sql(db *xorm.Engine) {
	fmt.Println("--------------------query sql ----------------------")
	sql := "select * from goods"
	results, err := db.Query(sql)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, v := range results {
			fmt.Printf("good id:%v;\nprice:%v\n", string(v["id"]), string(v["price"]))
		}
	}
}

func exec_sql(db *xorm.Engine) {
	fmt.Println("--------------------exec sql ----------------------")
	sql := "update goods set price=? where id=?"
	res, err := db.Exec(sql, 40, 1)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("res:", res)
	}
}

func delete_row(db *xorm.Engine) {
	fmt.Println("--------------------delete row ----------------------")
	affected, err := db.Where("id >?", 2).Delete(new(Goods))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("delete affected:", affected)
	}
}
func simple_transact(db *xorm.Engine) {
	fmt.Println("--------------------simple transact---------------------")
	session := db.NewSession()
	defer session.Close()
	// add Begin() before any action
	err := session.Begin()
	good := Goods{Id: 5, Price: 70}
	_, err = session.Insert(&good)
	if err != nil {
		session.Rollback()
		fmt.Println("error:", err)
		return
	}
	good2 := Goods{Price: 40}
	_, err = session.Where("id = ?", 2).Update(&good2)
	if err != nil {
		session.Rollback()
		fmt.Println("error:", err)
		return
	}

	_, err = session.Exec("delete from goods where id = ?", 2)
	if err != nil {
		session.Rollback()
		fmt.Println("error:", err)
		return
	}

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

func op_table_test() {
	db, err := xorm.NewEngine("postgres", "postgres://postgres:postgres@localhost:5432/test?sslmode=disable")
	if err != nil {
		fmt.Printf("%v \n", err)
	} else {
		defer db.Close()
		var ok bool
		ok, err = db.IsTableExist("goods")
		if err != nil {
			fmt.Println("table goods IsTableExist error", err)
		} else {
			fmt.Println("table goods is exist :", ok)
			err = db.Sync2(&Goods{})
			if err != nil {
				fmt.Println("table goods Sync2 error", err)
			} else {
				ok, err = db.IsTableExist("goods")
				if err != nil {
					fmt.Println("table goods IsTableExist error", err)
				} else {
					fmt.Println("table goods is exist :", ok)
					insert_table(db)
					update_table(db)
					select_get_table(db)
					select_exist_table(db)
					select_find_table(db)
					select_iterate_table(db)
					select_row_table(db)
					sum_table(db)
					count_table(db)
					exec_sql(db)
					query_sql(db)
					simple_transact(db)
					delete_row(db)
					err = db.DropTables("goods")
					if err != nil {
						fmt.Println("table goods drop error:", err)
					} else {
						fmt.Println("table drop ok")
						ok, err = db.IsTableExist("goods")
						if err != nil {
							fmt.Println("table goods IsTableExist error", err)
						} else {
							fmt.Println("table goods is exist :", ok)
						}
					}
				}
			}
		}
	}
}
func main() {
	sync_table()
	op_table_test()
}
