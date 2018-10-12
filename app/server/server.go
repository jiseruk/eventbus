package server

func Init() {
	r := GetRouter()
	r.Run(":8080")
}