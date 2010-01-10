// vim: ft=go

package main

func sum(a []int) int {
   s := 0
   for i := 0; i < len(a); i++ {
      s += a[i]
   }
   return s
}

func main() {
   s := sum([]int{1,2,3})
   println(s)
}
