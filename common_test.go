package gutil

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDataSlice(t *testing.T) {
	strs := []string{"a", "b", "c", "d", "e"}

	slice := DataSlice(strs, 2)
	fmt.Println(slice)
}

func TestMd5(t *testing.T) {
	slice := Md5("abcd")
	fmt.Println(slice)
}

func TestJJ(t *testing.T) {
	for i := 0; i < 10; i++ {
		problem := generateMathProblem(100)
		fmt.Println(problem)
	}
}

// 生成随机整数
func generateRandomNumber(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

// 生成随机运算符
func generateRandomOperator() string {
	operators := []string{"+", "-", "*", "/"}
	return operators[rand.Intn(len(operators))]
}

// 生成题目
func generateMathProblem(max int) string {
	num1 := generateRandomNumber(max)
	num2 := generateRandomNumber(max)
	operator := generateRandomOperator()

	// 防止除法出现小数，并且确保除数不为0
	if operator == "/" {
		for num2 == 0 || num1%num2 != 0 {
			num1 = generateRandomNumber(max)
			num2 = generateRandomNumber(max)
		}
	}

	return fmt.Sprintf("%d %s %d =", num1, operator, num2)
}
