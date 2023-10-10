/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package iris

type TextBlocks map[string]TextBlock

func UserInterfaceStartThread(ch_data_input chan MessageAndCharactersForWindowsUpdate, dataInputLineSeparator string) {
	userInterfaceInit()
	// ch_user_input := make(chan string)
	// ch_terminal_size_change_detect := make(chan [2]int)

	textBlocks := TextBlocks{}
	go channelRead_dataInputInterpret(ch_data_input, &textBlocks, dataInputLineSeparator)

}

func userInterfaceInit() {
	terminal_console_clear()
	terminal_console_input_buffering_disable()
	terminal_console_character_hide()
}

func UserInterfaceExitThread() {
	terminal_console_character_show()
	terminal_console_input_buffering_enable()
}
