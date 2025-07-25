package main

import (
	"fmt"
	"ai-context-cli/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("📂 === ADD CONTEXT TO ALL FILES === 📂")
	fmt.Println("✅ Escaneado completo de proyecto")
	fmt.Println("✅ Progreso en tiempo real")
	fmt.Println("✅ Conteo de archivos y tiempo estimado")
	fmt.Println("✅ Exclusión de tipos de archivo")
	fmt.Println("✅ Generación de contexto comprehensivo")
	fmt.Println()
	
	fmt.Println("🎯 FUNCIONALIDADES:")
	fmt.Println("🔍 Escaneo inteligente de archivos")
	fmt.Println("📊 Estadísticas en tiempo real")
	fmt.Println("🚫 Exclusión automática de archivos innecesarios")
	fmt.Println("📝 Generación de contexto estructurado")
	fmt.Println("🎨 UI profesional con progreso visual")
	fmt.Println()
	
	fmt.Println("📋 LO QUE VERÁS:")
	fmt.Println("1. Navegación a 'Add Context to All Files'")
	fmt.Println("2. Spinner inicial y estimación de archivos")
	fmt.Println("3. Barra de progreso durante el escaneo")
	fmt.Println("4. Generación de contexto comprehensive")
	fmt.Println("5. Resultado final con estadísticas")
	fmt.Println()
	
	fmt.Println("🎮 INSTRUCCIONES:")
	fmt.Println("1. Selecciona '📂 Add Context to All Files'")
	fmt.Println("2. Observa el progreso del escaneo")
	fmt.Println("3. Ve la generación de contexto")
	fmt.Println("4. Revisa el resultado final")
	fmt.Println("5. Usa ESC para regresar al menú")
	fmt.Println()
	
	fmt.Println("⚠️  NOTA: El escaneo será del directorio actual")
	fmt.Println("🚀 Empezando la aplicación...")
	fmt.Println()
	
	model := app.NewModel()
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}