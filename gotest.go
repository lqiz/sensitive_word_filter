package main

func main() {
	s := "ANC"
	b := []byte(s)
	s2 := string(b)
	for _, v := range s2 {
		print( v )
		//do something with i,v
	}
}
