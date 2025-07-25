package main

import (
	"fmt"
	"ai-context-cli/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("🧭 === SISTEMA DE NAVEGACIÓN === 🧭")
	fmt.Println("✅ Breadcrumb trail en la parte superior")
	fmt.Println("✅ Funcionalidad de back con ESC")
	fmt.Println("✅ Stack de historial de navegación")
	fmt.Println("✅ Indicadores visuales de navegación")
	fmt.Println("✅ Patrones consistentes")
	fmt.Println()
	
	fmt.Println("🎯 FUNCIONALIDADES:")
	fmt.Println("🏠 Main Menu → Contexto Engine (sin back)")
	fmt.Println("📂 Add Context All → Contexto Engine › Add Context › All Files")
	fmt.Println("📁 Add Context Folder → Contexto Engine › Add Context › Folder") 
	fmt.Println("📋 Context Preview → Contexto Engine › Context Preview")
	fmt.Println("🤖 Model Selection → Contexto Engine › Model Selection")
	fmt.Println()
	
	fmt.Println("🔄 NAVEGACIÓN:")
	fmt.Println("- Breadcrumbs muestran la ruta actual")
	fmt.Println("- ESC regresa a pantalla anterior")
	fmt.Println("- '← ESC: Back' aparece cuando hay historial")
	fmt.Println("- Stack mantiene historial completo")
	fmt.Println("- Navegación centrada y profesional")
	fmt.Println()
	
	fmt.Println("🎮 INSTRUCCIONES:")
	fmt.Println("1. Selecciona cualquier opción del menú")
	fmt.Println("2. Observa los breadcrumbs en la parte superior")
	fmt.Println("3. Presiona ESC para regresar al menú")
	fmt.Println("4. Navega entre diferentes secciones")
	fmt.Println()
	
	model := app.NewModel()
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}