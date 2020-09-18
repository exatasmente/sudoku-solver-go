package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"os"
	"bufio"
	"math/rand"
)
type DPLL struct{	
	literais []int
	resultado []Clausula
	formula chan []Clausula
}
type Clausula struct {
	valor []int
}
type Sudoku struct {
	valor  [9][9]int
	valido bool
}

type Regras struct {
	clausulas []Clausula
}

func (dpll *DPLL) Dpll(formula []Clausula) (bool) {	
	index := rand.Intn(len(formula))
	el := rand.Intn(len(formula[index].valor))
	var elemento int = formula[index].valor[el]
	var formulaAux = dpll.Simplifica(formula,elemento)
	dpll.resultado = formulaAux

	
	if len(formulaAux) == 0 {
		return true
	}else{
		
		for i:= 0; i<len(formulaAux) ; i++{
			var count = 0
			fmt.Println(formulaAux[i].valor)
			for j:= 0; j<len(formulaAux[i].valor) ; j++{
				
				if formulaAux[i].valor[j] != 0{
					count++
				}
			}
			if count == 0 {
				return false
			}
		}
		
	}
	
	
	var clausula Clausula
	clausula.valor = append(clausula.valor,elemento)
	if  dpll.Dpll(append(formulaAux,clausula)) == true {
		
		return true
	}else{
		clausula.valor[0] = clausula.valor[0]*(-1)
		if dpll.Dpll(append(formulaAux,clausula)) {
			return true
		}else{
			return false
		}
	}

}

func (dpll *DPLL)  Init(){
	dpll.literais = literais()
	
}


func literais() []int {
	var lista []int
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			for z := 1; z < 10; z++ {
				lista = append(lista, concatena(x,y,z, true)) // so está adicionando literais positivos tem que adicionar os negativos também se forem Valor_Verdade
				lista = append(lista, concatena(x,y,z, false))
			}			
		}
	}
	return lista
}

 func separaLiterais(ch chan [9][9]int){
 	var sudoku [9][9]int
	linhas,err := lerArquivo("sudoku-logica/facil.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for b:=0 ; b<9; b++{
		for a:=0 ; a<9; a++{
			sudoku[b][a],_= strconv.Atoi(linhas[a])
		}
	}
	ch <- sudoku


}

func(dpll *DPLL) Simplifica(formula []Clausula,elemento int) []Clausula {
	for i := 0 ; i < len(formula) ; i++ {
		for j := 0 ; j < len(formula[i].valor) ; j++ {
			if formula[i].valor[j] == elemento {
				formula = append(formula[:i],formula[i+1:]...) 
				break
			}else if formula[i].valor[j] == elemento*(-1) {
				if len(formula[i].valor) > 1 {
					formula[i].valor[j]  = 0
				}else{
					formula = append(formula[:i],formula[i+1:]...) 
					break
				}

			}
		}
	}
	return formula
}

func concatena(x,y,z int,positivo bool) int{

	if(!positivo){
		return  ((x*100) + (y*10) + (z) )*(-1)
	}
	return  (x*100) + (y*10) + (z)
}


func restricao1(ch chan Clausula, quit chan int) {
	
	for x := 1; x < 10; x++ {
	
		for y := 1; y < 10; y++ {
			var cl Clausula
			for z := 1; z < 10; z++ {
				aux := concatena(x, y, z, true)
				cl.valor = append(cl.valor,aux)
			}
			ch <-cl

		}

	}
	quit <- 1
}


func restricao2(ch chan Clausula, quit chan int) {
	
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			for z := 1; z < 9; z++ {
				
				for i := x + 1; i < 10; i++ {
					var cl Clausula
					aux := concatena(z , x , y, false)
					cl.valor = append(cl.valor,aux)
					aux = concatena(i , x , y, false)
					cl.valor = append(cl.valor,aux)
					ch <- cl
				}
				
			}
		}
	}
	quit <- 1
}

func restricao3(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			for z := 1; z < 9; z++ {
				
				for i := y + 1; i < 10; i++ {
					var cl Clausula
					aux:= concatena(x,z,y,false)
					cl.valor = append(cl.valor, aux)
					aux = concatena(x,i,y,false)
					cl.valor = append(cl.valor,aux)
					ch <- cl
				}
				
			}
		}
	}
	quit <- 1
}
func restricao4(ch chan Clausula, quit chan int) {

	for x := 1; x < 10; x++ {
		for y := 0; y < 3; y++ {
			for z := 0; z < 3; z++ {				
				for i := 1; i < 4; i++ {					
					for j := 1; j < 4; j++ {
						
						for a := j + 1; a < 4; a++ {
							var cl Clausula
							aux := concatena(3*y+i,3*z+j,x,false)
							cl.valor = append(cl.valor, aux)
							aux = concatena(3*y+i,3*z+a,x,false)
							cl.valor = append(cl.valor, aux)
							ch <- cl
						}
						
					}

				}
			}
		}
	}
	quit <- 1
}

func restricao5(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for y := 0; y < 3; y++ {
			for z := 0; z < 3; z++ {
				for i := 1; i < 4; i++ {
					for j := 1; j < 4; j++ {
						for a := i + 1; a < 4; a++ {
							
							for b := 1; b < 4; b++ {
								var cl Clausula
								aux := concatena(3*y+i, 3*z+j,x,false)
								cl.valor = append(cl.valor,aux)
								aux = concatena(3*y+a, 3*z+b, x, false)
								cl.valor = append(cl.valor,aux)
								ch <- cl
							}
							

						}
					}

				}
			}
		}
	}
	quit <- 1
}

func restricao6(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			for z := 1; z < 9; z++ {				
				for i := z + 1; i < 10; i++ {
					var cl Clausula
					aux := concatena(x, y, z,false)
					cl.valor = append(cl.valor,aux)
					aux = concatena(x, y, i, false)
					cl.valor = append(cl.valor,aux)
					ch <- cl					
				}
				
			}
		}
	}
	quit <- 1
}

func restricao7(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			var cl Clausula
			for z := 1; z < 10; z++ {
				aux := concatena(z, x, y,true)
				cl.valor = append(cl.valor,aux)
			}
			ch <- cl
		}
	}
	quit <- 1
}
func restricao8(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			var cl Clausula
			for z := 1; z < 10; z++ {
				aux := concatena(x, z, y,true)
				cl.valor = append(cl.valor,aux)
			}
			ch <- cl
		}
	}
	quit <- 1
}
func restricao9(ch chan Clausula, quit chan int) {

	for x := 1; x < 10; x++ {
		for y := 0; y < 3; y++ {
			for z := 0; z < 3; z++ {
				var cl Clausula
				for i := 1; i < 4; i++ {					
					for j := 1; j < 4; j++ {
						aux := concatena(3*y+i, 3*z+j, x,true)
						cl.valor = append(cl.valor,aux)
					}					
				}
				ch <- cl

			}
		}
	}
	quit <- 1
}
func geraClausula(regra int, cb chan Clausula, wg *sync.WaitGroup) {

	
	quit := make(chan int)
	clChan := make(chan Clausula)
	
	switch regra {
	 case 1:
	 	go restricao1(clChan, quit)
		
	case 2:
	 	go restricao2(clChan, quit)
	case 3:
		go restricao3(clChan, quit)
	case 4:
		go restricao4(clChan, quit)
	case 5:
		go restricao5(clChan, quit)
	case 6:
		go restricao6(clChan, quit)
	case 7:
		go restricao7(clChan, quit)
	case 8:
		go restricao8(clChan, quit)
	case 9:
		go restricao9(clChan, quit)
	default:
		fmt.Println("inválido")
		return
	}
	for {
		select {
		case <-quit:
			wg.Done()
			return
		case cl := <-clChan:
			cb <- cl
		}
	}


}
func lerArquivo(caminhoDoArquivo string) ([]string, error) {
	
	arquivo, err := os.Open(caminhoDoArquivo)
	
	if err != nil {
		return nil, err
	}
	defer arquivo.Close()
	var linhas []string
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		linhas = append(linhas, scanner.Text())
	}

	// Retorna as linhas lidas e um erro se ocorrer algum erro no scanner
	return linhas, scanner.Err()
}




func main() {
	
	var wg sync.WaitGroup
	
    
	// cn := make (chan[9][9]int)
	// wg.Add(1)
	// go separaLiterais(cn)
	// go func() {
	// 	for{
	// 		select {
	// 		case cl := <-cn:
	// 			fmt.Println(cl)
	// 			wg.Done()
	// 			return
	// 		}

	// 	}
	// }()
	
	// wg.Wait()
	
	
	var s Sudoku
	s.valido = false
	var r Regras
	cb := make(chan Clausula)
	done := make(chan int)
	t1 := time.Now()
	 for i := 1; i <10; i++ {
		wg.Add(1)
		go geraClausula(i, cb, &wg)

	}
	running := true
	go func() {
		for running {
			select {
			case cl := <-cb:
				
				r.clausulas = append(r.clausulas,cl)
			case <-done:
				fmt.Printf("Done \n")
				return
			default:

			}

		}
	}()
	wg.Wait()
	done <- 1

	var dpll DPLL
	dpll.Init()

	fmt.Println(dpll.Dpll(r.clausulas))
	
	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Println(diff)
		
	

}