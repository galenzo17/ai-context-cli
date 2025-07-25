package main

import (
	"fmt"
	"ai-context-cli/internal/app"
)

func main() {
	fmt.Println("🚀 === MENÚ PROFESIONAL MEJORADO === 🚀")
	fmt.Println("✅ Iconos descriptivos para cada opción")
	fmt.Println("✅ Textos claros y profesionales:")
	fmt.Println("   📂 Add Context to All Files - Scan entire project")
	fmt.Println("   📁 Add Context to Specific Folder - Choose folder")
	fmt.Println("   📋 Preview Context Before Sending - Review context")
	fmt.Println("   🤖 Select AI Model - Configure models")
	fmt.Println("   🚪 Exit - Quit application")
	fmt.Println("✅ Sistema de ayuda con '?' - Modal con detalles")
	fmt.Println("✅ Navegación intuitiva y clara")
	fmt.Println("✅ Mucho más contexto para el usuario")
	fmt.Println()
	
	model := app.NewModel()
	
	fmt.Println("🎯 RESULTADO FINAL:")
	view := model.View()
	fmt.Print(view)
	
	fmt.Println("\n🏆 CARACTERÍSTICAS PROFESIONALES:")
	fmt.Println("- Cada botón explica claramente qué hace")
	fmt.Println("- Iconos ayudan a identificar rápidamente las opciones")
	fmt.Println("- Sistema de ayuda detallado con '?' key")
	fmt.Println("- Textos descriptivos y profesionales")
	fmt.Println("- Banner compacto pero visible")
	fmt.Println("- Navegación intuitiva")
	fmt.Println()
	fmt.Println("¡MUCHO más profesional y fácil de usar! 🌟")
}