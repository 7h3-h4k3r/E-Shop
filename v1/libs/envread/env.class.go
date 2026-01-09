package envread   

import (
	"encoding/json"
	"os"
)


type venv struct{
	MongoDB string `json:"MONGO_DB"`
	Collection string `json:"COLL_DB"`
	DataBase string `json:"DB"`
	JwtKey string `json:"SEC_KEY"`
}


var  Env_file_read venv;

func Getenv(){
	
	file , err := os.ReadFile("../venv.json")
	
	if err!=nil{
		panic(err)
	}

	
	if err := json.Unmarshal(file,&Env_file_read);err!=nil{
		panic(err)
	}

}