package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Literal struct {
	elemento int
}
type Clausula struct {
	literais []Literal
}
type Formula struct {
	clausulas  []Clausula
	literais   []int
	valoracao  map[int]int
	polaridade map[int]int
	frequencia map[int]int
}

type Sudoku struct {
	valor  [9][9]int
	valido bool
}
type DPLL struct {
	resultado   []Literal
	literaisQty int
	clasulasQty int
}

func printFormula(resultado Formula, sudoku [9][9]int) {
	var lista []int
	for k, val := range resultado.valoracao {
		if val == 1 {
			lista = append(lista, k)
		}
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudoku[i][j] > 0 {
				lista = append(lista, concatena(i+1, j+1, sudoku[i][j], true).elemento)
			}
		}
	}
	sort.Ints(lista)
	index := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%d ", (lista[index]%10)%10)
			index++
		}
		fmt.Println()
	}
}
func (dpll *DPLL) Init(sudoku [9][9]int, formula *Formula) {
	formula.valoracao = make(map[int]int)
	formula.frequencia = make(map[int]int)
	formula.polaridade = make(map[int]int)
	var valoracoes []Literal
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudoku[i][j] != 0 {
				var literal = concatena(i+1, j+1, sudoku[i][j], true)
				valoracao := literal
				valoracoes = append(valoracoes, valoracao)
			}
		}
	}

	for i := 0; i < len(valoracoes); i++ {
		formula.clausulas = append(formula.clausulas, Clausula{[]Literal{valoracoes[i]}})
		formula.valoracao[valoracoes[i].elemento] = 1
	}
	for _, c := range formula.clausulas {
		for _, l := range c.literais {
			if valor, ok := formula.frequencia[l.elemento]; !ok {
				formula.frequencia[l.elemento]++
				formula.valoracao[l.elemento] = -1
				if l.elemento > 0 {
					formula.polaridade[l.elemento]++
				} else {
					formula.polaridade[l.elemento]--
				}
			} else {
				valor++
				if l.elemento > 0 {
					formula.polaridade[l.elemento]++
				} else {
					formula.polaridade[l.elemento]--
				}
			}

		}
	}

}

func (dpll *DPLL) _UnitPropagate(formula Formula) (int, Formula) {
	var loop = true
	if len(formula.clausulas) == 0 {

		return 1, formula
	}

	for loop {
		loop = false
		for i := 0; i < len(formula.clausulas); i++ {
			if len(formula.clausulas[i].literais) == 1 && formula.clausulas[i].literais[0].elemento > 1 {
				var valid = -1
				loop = true
				literal := formula.clausulas[i].literais[0]
				formula.frequencia[formula.clausulas[i].literais[0].elemento] = -1

				valid, formula = dpll._ApplyTransform(formula, literal.elemento)
				if valid != -1 {
					return valid, formula
				}
				break
			} else if len(formula.clausulas) == 0 {
				return 0, formula
			}
		}
	}
	return -1, formula
}

func (dpll *DPLL) _ApplyTransform(formula Formula, literal int) (int, Formula) {
	elemento := literal
	for i := 0; i < len(formula.clausulas); i++ {
		for j := 0; j < len(formula.clausulas[i].literais); j++ {
			if formula.clausulas[i].literais[j].elemento == elemento {
				formula.clausulas[i] = Clausula{[]Literal{Literal{1}}}
			} else if formula.clausulas[i].literais[j].elemento == elemento*(-1) {
				formula.clausulas[i].literais[j] = Literal{0}

			}

		}

	}
	return -1, formula
}
func (dpll *DPLL) _PegaLiteral(formula Formula) int {
	var max = 0
	var c = -1
	var l = 0 - 1
	for i := 0; i < len(formula.clausulas); i++ {
		for j := 0; j < len(formula.clausulas[i].literais); j++ {
			if formula.frequencia[formula.clausulas[i].literais[j].elemento] > max && formula.valoracao[formula.clausulas[i].literais[j].elemento] == -1 {
				max = formula.frequencia[formula.clausulas[i].literais[j].elemento]
				c, l = i, j
			}

		}
	}
	if c != -1 {
		return formula.clausulas[c].literais[l].elemento
	}
	return c
}
func (dpll *DPLL) Dpll(formula Formula) (int, Formula) {

	result, formula := dpll._UnitPropagate(formula)
	if result == 1 {
		return 1, formula
	} else if result == 0 {
		return 0, formula
	}

	literal := dpll._PegaLiteral(formula)
	if literal == -1 {
		return 0, formula
	}

	for j := 0; j < 2; j++ {
		var novaFormula Formula = formula
		if novaFormula.polaridade[literal] > 0 {
			novaFormula.valoracao[literal] = 1
		} else {
			novaFormula.valoracao[literal*-1] = 0
		}
		novaFormula.frequencia[literal] = -1
		result, novaFormula := dpll._ApplyTransform(novaFormula, literal)
		if result == 1 {
			return 1, novaFormula
		} else if result == 0 {
			continue
		}
		dpllResult, novaFormula := dpll.Dpll(novaFormula)
		if dpllResult == 1 {
			return dpllResult, novaFormula
		}
	}
	return -1, formula
}

func literais() []Literal {
	var lista []Literal
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			for z := 1; z < 10; z++ {
				lista = append(lista, concatena(x, y, z, true))
				lista = append(lista, concatena(x, y, z, false))
			}
		}
	}
	return lista
}

func concatena(x, y, z int, positivo bool) Literal {

	if !positivo {
		return Literal{((x * 100) + (y * 10) + (z)) * (-1)}
	}
	return Literal{(x * 100) + (y * 10) + (z)}
}

func restricao1(ch chan Clausula, quit chan int) {

	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			var cl Clausula
			for z := 1; z < 10; z++ {
				aux := concatena(x, y, z, true)
				cl.literais = append(cl.literais, aux)
			}
			ch <- cl

		}

	}
	quit <- 1
}

func restricao2(ch chan Clausula, quit chan int) {

	for y := 1; y < 10; y++ {
		for z := 1; z < 10; z++ {
			for x := 1; x < 9; x++ {
				for i := x + 1; i < 10; i++ {
					var cl Clausula
					aux := concatena(x, y, z, false)
					cl.literais = append(cl.literais, aux)
					aux = concatena(i, y, z, false)
					cl.literais = append(cl.literais, aux)
					ch <- cl
				}

			}
		}
	}
	quit <- 1
}

func restricao3(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for z := 1; z < 10; z++ {
			for y := 1; y < 9; y++ {
				for i := y + 1; i < 10; i++ {
					var cl Clausula
					aux := concatena(x, z, y, false)
					cl.literais = append(cl.literais, aux)
					aux = concatena(x, i, y, false)
					cl.literais = append(cl.literais, aux)
					ch <- cl
				}

			}
		}
	}
	quit <- 1
}
func restricao4(ch chan Clausula, quit chan int) {

	for z := 1; z < 10; z++ {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for x := 1; x < 4; x++ {
					for y := 1; y < 4; y++ {
						for k := y + 1; k < 4; k++ {
							var cl Clausula
							aux := concatena(3*i+x, 3*j+y, z, false)
							cl.literais = append(cl.literais, aux)
							aux = concatena(3*i+x, 3*j+k, z, false)
							cl.literais = append(cl.literais, aux)
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
	for z := 1; z < 10; z++ {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for x := 1; x < 4; x++ {
					for y := 1; y < 4; y++ {
						for k := i + 1; k < 4; k++ {
							for l := 1; l < 4; l++ {
								var cl Clausula
								aux := concatena(3*i+x, 3*j+y, z, false)
								cl.literais = append(cl.literais, aux)
								aux = concatena(3*i+k, 3*j+l, z, false)
								cl.literais = append(cl.literais, aux)
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
					aux := concatena(x, y, z, false)
					cl.literais = append(cl.literais, aux)
					aux = concatena(x, y, i, false)
					cl.literais = append(cl.literais, aux)
					ch <- cl
				}

			}
		}
	}
	quit <- 1
}

func restricao7(ch chan Clausula, quit chan int) {
	for y := 1; y < 10; y++ {
		for z := 1; z < 10; z++ {
			var cl Clausula
			for x := 1; x < 10; x++ {
				aux := concatena(x, y, z, true)
				cl.literais = append(cl.literais, aux)
			}
			ch <- cl
		}
	}
	quit <- 1
}
func restricao8(ch chan Clausula, quit chan int) {
	for x := 1; x < 10; x++ {
		for z := 1; z < 10; z++ {
			var cl Clausula
			for y := 1; y < 10; y++ {
				aux := concatena(x, y, z, true)
				cl.literais = append(cl.literais, aux)
			}
			ch <- cl
		}
	}
	quit <- 1
}
func restricao9(ch chan Clausula, quit chan int) {

	for z := 1; z < 10; z++ {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				var cl Clausula
				for x := 1; x < 4; x++ {
					for y := 1; y < 4; y++ {
						aux := concatena(3*i+x, 3*j+y, z, true)
						cl.literais = append(cl.literais, aux)
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
func lerArquivo(caminhoDoArquivo string) ([][]string, error) {

	arquivo, err := os.Open(caminhoDoArquivo)

	if err != nil {
		return nil, err
	}
	defer arquivo.Close()
	var linhas [][]string
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		linhas = append(linhas, strings.Split(scanner.Text(), " "))
	}

	// Retorna as linhas lidas e um erro se ocorrer algum erro no scanner
	return linhas, scanner.Err()
}

func separaLiterais(ch chan [9][9]int) {
	var sudoku [9][9]int
	linhas, err := lerArquivo("sudoku-logica/dificil.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for b := 0; b < len(linhas); b++ {
		for a := 0; a < len(linhas[b]); a++ {
			sudoku[b][a], _ = strconv.Atoi(linhas[b][a])
		}
	}
	ch <- sudoku

}

func main() {

	var wg sync.WaitGroup
	var dpll DPLL
	var sudoku [9][9]int
	cn := make(chan [9][9]int)

	wg.Add(1)
	go separaLiterais(cn)
	go func() {
		for {
			select {
			case sudoku = <-cn:
				wg.Done()
				return
			}

		}
	}()

	wg.Wait()
	var formula Formula
	cb := make(chan Clausula)
	done := make(chan int)

	running := true
	go func() {
		for running {
			select {
			case cl := <-cb:
				formula.clausulas = append(formula.clausulas, cl)
			case <-done:
				fmt.Printf("Clausulas Geradas \n")
				return
			default:

			}

		}
	}()

	for i := 1; i < 6; i++ {
		wg.Add(1)
		go geraClausula(i, cb, &wg)
		//wg.Wait() // gera clasulas assincrono
	}

	wg.Wait()
	done <- 1
	dpll.Init(sudoku, &formula)

	t1 := time.Now()
	solved, r := dpll.Dpll(formula)
	t2 := time.Now()
	if solved == 1 {
		fmt.Println("SAT")
	} else {
		fmt.Println("UNSAT")
	}
	printFormula(r, sudoku)
	diff := t2.Sub(t1)
	fmt.Println("Executado em :", diff)

}
