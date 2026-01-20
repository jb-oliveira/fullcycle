package main

import (
	"context"
	"fmt"
	"time"
)

// longRunningTask simula uma tarefa que executa várias etapas (iterações).
// Ela verifica o contexto a cada passo para ver se deve parar.
func longRunningTask(ctx context.Context, taskName string, maxIterations int) (string, error) {
	fmt.Printf("[%s] Tarefa iniciada. Vai tentar %d iterações.\n", taskName, maxIterations)

	for i := 1; i <= maxIterations; i++ {
		// 1. O Ponto de Verificação Crítico:
		select {
		case <-ctx.Done():
			// O canal ctx.Done() foi fechado, o que significa que o timeout expirou
			// ou o contexto foi cancelado.
			fmt.Printf("[%s] ❌ Interrompida! Contexto cancelado após %d iterações.\n", taskName, i-1)
			return "", ctx.Err() // Retorna o erro do contexto (deadline exceeded)
		default:
			// Não houve cancelamento, continua o trabalho.
		}

		// Simulação de uma parte do trabalho (cada passo demora 100ms)
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[%s] Passo %d/%d concluído.\n", taskName, i, maxIterations)
	}

	// Se o loop terminar sem o contexto ser cancelado
	fmt.Printf("[%s] ✅ Concluída com sucesso após %d iterações.\n", taskName, maxIterations)
	return fmt.Sprintf("Tarefa '%s' concluída.", taskName), nil
}

func main() {

	// --- Cenário 1: Sucesso (Timeout longo o suficiente) ---
	fmt.Println("--- Cenário 1: Sucesso (Timeout de 1 segundo) ---")

	// Definimos 1 segundo de timeout. A tarefa completa (5 passos * 100ms = 500ms).
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()

	result1, err1 := longRunningTask(ctx1, "Task-A", 5)

	if err1 != nil {
		fmt.Printf("Resultado Task-A: Erro: %v\n", err1)
	} else {
		fmt.Printf("Resultado Task-A: %s\n", result1)
	}

	fmt.Println("\n" + "--- Cenário 2: Timeout Expirado (Timeout curto) ---")

	// --- Cenário 2: Timeout Expirado ---

	// Definimos apenas 250ms de timeout.
	// A tarefa precisaria de 1 segundo (10 passos * 100ms).
	ctx2, cancel2 := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel2()

	result2, err2 := longRunningTask(ctx2, "Task-B", 10)

	if err2 != nil {
		fmt.Printf("Resultado Task-B: Erro: %v\n", err2)
	} else {
		fmt.Printf("Resultado Task-B: %s\n", result2)
	}
}
