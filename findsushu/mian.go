package main

//generate 2,3,4,...to chan in
func Generate(in chan int) {
	for num := 2; ; num++ {
		in <- num
	}
}

//fiter prime number
func Fiter(in chan int, out chan int, prime int) {
	for {
		temp := <-in
		if temp%prime != 0 {
			out <- temp
		}
	}
}

func main() {
	in := make(chan int)
	go Generate(in)
	for i := 0; i < 1000; i++ {
		prime := <-in
		println(prime)
		out := make(chan int)
		go Fiter(in, out, prime)
		println(out)
		in = out
	}
}
