package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"github.com/jcelliott/lumber"
)

const Version = "1.0"

type (
	Logger interface{
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}


	Driver struct{
		Mutex *sync.Mutex
		Mutexes map[string] *sync.Mutex
		Dir string
		Log Logger
	}

	Option struct{
		Logger
	}
)

func New(dir string, options *Option) (*Driver, error){
	dir = filepath.Clean(dir)

	opts := Option{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		Dir: dir,
		Mutexes: make(map[string] *sync.Mutex),
		Log: opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil{
		opts.Logger.Debug("Using '%s' (Database already exists...)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating database at %s ...", dir)
	return &driver, os.MkdirAll(dir, 0755)
}


func (d *Driver)Write(collection string, resource string, v interface{}) error{

	if collection == nil {
		fmt.Errorf("No collection specified...")
	}

	if resourse == nil {
		fmt.Errorf("No resource name specified...")
	}

	mutex := d.getOrCreateMutex(collection)

	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.Dir, collection)
	fnlPath := filepath.Join(dir, resource + ".json")
	tmpPath := filepath.Join(fnlPath + ".tmp")

	if err := os.MkdirAll(dir, 0755); err != nil{
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil{
		return err
	}

	b = append(b, "\n")

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil{
		return err
	}

	return os.Rename(tmpPath, fnlPath)

}



func stat(path string) (fi *os.FileInfo, err error){
	if fi, err = os.Stat(path); os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return
}
  

type Address struct{
	City string
	State string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {

	dir := "./"

	db, err := New(dir, nil)
	if err != nil{
		fmt.Println("Err", err)
	}

	employees := []User{
		{"John", "23", "23344333", "Myrl Tech", Address{"bangalore", "karnataka", "india", "410013"}},
		{"Paul", "25", "23344333", "Google", Address{"san francisco", "california", "USA", "410013"}},
		{"Robert", "27", "23344333", "Microsoft", Address{"bangalore", "karnataka", "india", "410013"}},
		{"Vince", "29", "23344333", "Facebook", Address{"bangalore", "karnataka", "india", "410013"}},
		{"Neo", "31", "23344333", "Remote-Teams", Address{"bangalore", "karnataka", "india", "410013"}},
		{"Albert", "32", "23344333", "Dominate", Address{"bangalore", "karnataka", "india", "410013"}},
	}

	for _, user := range employees{
		db.Write("user", user.Name, User{
			Name: user.Name,
			Age: user.Age,
			Contact: user.Contact,
			Company: user.Company,
			Address: user.Address,
		})
	}

	records, err := db.ReadAll("user")
	if err != nil{
		fmt.Println("Err", err)
	}

	users := []User{}

	for _, record := range records{
		employees := User{}

		if err := json.Unmarshal([]byte(record), &employees); err != nil{
			fmt.Println("err", err)
		}

		users = append(users, employees)


	}


}
