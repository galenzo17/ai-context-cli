package main

import (
	"fmt"
	"ai-context-cli/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("🚀 === SISTEMA DE CARGA Y FEEDBACK === 🚀")
	fmt.Println("✅ Componentes de spinner para estados de carga")
	fmt.Println("✅ Barras de progreso para operaciones de archivos")
	fmt.Println("✅ Sistema de mensajes de éxito/error")
	fmt.Println("✅ Notificaciones tipo toast con auto-dismiss")
	fmt.Println("✅ Delays de simulación para mostrar estados")
	fmt.Println("✅ Feedback visual profesional")
	fmt.Println()
	
	fmt.Println("🎯 FUNCIONALIDADES:")
	fmt.Println("📂 Add Context to All Files - Muestra progress bar")
	fmt.Println("📁 Add Context to Folder - Simulación de proceso")
	fmt.Println("📋 Context Before - Loading spinner")
	fmt.Println("🤖 Select Model - Toast notifications")
	fmt.Println()
	
	fmt.Println("🔄 INTERACCIÓN:")
	fmt.Println("- Selecciona cualquier opción para ver loading")
	fmt.Println("- Observa spinners, progress bars y toasts")
	fmt.Println("- Todo regresa al menú automáticamente")
	fmt.Println()
	
	model := app.NewModel()
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}