package main

import (
	"fmt"
	"strconv"
)

var ppIcons []string

//Generate the icons arrays and assign it to ppicons
func IconsArray() {
	arr := []string{"https://image.flaticon.com/icons/png/512/4652/4652622.png", "https://image.flaticon.com/icons/png/512/4652/4652623.png", "https://image.flaticon.com/icons/png/512/4652/4652624.png", "https://image.flaticon.com/icons/png/512/4652/4652625.png", "https://image.flaticon.com/icons/png/512/4652/4652626.png", "https://image.flaticon.com/icons/png/512/4652/4652627.png", "https://image.flaticon.com/icons/png/512/4652/4652628.png", "https://image.flaticon.com/icons/png/512/4652/4652629.png", "https://image.flaticon.com/icons/png/512/4652/4652630.png", "https://image.flaticon.com/icons/png/512/4652/4652631.png", "https://image.flaticon.com/icons/png/512/4652/4652632.png", "https://image.flaticon.com/icons/png/512/4652/4652633.png", "https://image.flaticon.com/icons/png/512/4652/4652634.png", "https://image.flaticon.com/icons/png/512/4652/4652635.png", "https://image.flaticon.com/icons/png/512/4652/4652636.png", "https://image.flaticon.com/icons/png/512/4652/4652637.png", "https://image.flaticon.com/icons/png/512/4652/4652638.png", "https://image.flaticon.com/icons/png/512/4652/4652639.png", "https://image.flaticon.com/icons/png/512/4652/4652640.png", "https://image.flaticon.com/icons/png/512/4652/4652641.png", "https://image.flaticon.com/icons/png/512/4652/4652642.png", "https://image.flaticon.com/icons/png/512/4652/4652643.png", "https://image.flaticon.com/icons/png/512/4652/4652644.png", "https://image.flaticon.com/icons/png/512/4652/4652645.png", "https://image.flaticon.com/icons/png/512/4652/4652646.png", "https://image.flaticon.com/icons/png/512/4652/4652647.png", "https://image.flaticon.com/icons/png/512/4652/4652648.png", "https://image.flaticon.com/icons/png/512/4652/4652649.png", "https://image.flaticon.com/icons/png/512/4652/4652650.png", "https://image.flaticon.com/icons/png/512/4652/4652651.png", "https://image.flaticon.com/icons/png/512/4652/4652652.png", "https://image.flaticon.com/icons/png/512/4652/4652653.png", "https://image.flaticon.com/icons/png/512/4652/4652654.png", "https://image.flaticon.com/icons/png/512/4652/4652655.png", "https://image.flaticon.com/icons/png/512/4652/4652656.png", "https://image.flaticon.com/icons/png/512/4652/4652657.png", "https://image.flaticon.com/icons/png/512/4652/4652658.png", "https://image.flaticon.com/icons/png/512/4652/4652659.png", "https://image.flaticon.com/icons/png/512/4652/4652660.png", "https://image.flaticon.com/icons/png/512/4652/4652661.png", "https://image.flaticon.com/icons/png/512/4652/4652662.png", "https://image.flaticon.com/icons/png/512/4652/4652663.png", "https://image.flaticon.com/icons/png/512/4652/4652664.png", "https://image.flaticon.com/icons/png/512/4652/4652665.png", "https://image.flaticon.com/icons/png/512/4652/4652666.png", "https://image.flaticon.com/icons/png/512/4652/4652667.png", "https://image.flaticon.com/icons/png/512/4652/4652668.png", "https://image.flaticon.com/icons/png/512/4652/4652669.png", "https://image.flaticon.com/icons/png/512/4652/4652670.png", "https://image.flaticon.com/icons/png/512/4652/4652671.png"}
	ppIcons = arr
}

// return the src of the avatar img assigned
// to the profile picture of the user.
func NumberToPpIcon(icons string) string {
	avatarnb, err := strconv.Atoi(icons)
	if err != nil {
		fmt.Println(err)
	}
	result := ppIcons[avatarnb]
	return result
}
