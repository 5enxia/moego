package moego

import (
    "time"
)

func (e *Editor) PollTimerEvent() {
    for {
        switch <- e.timeChan {
        case RESET_MESSAGE:
            t := time.NewTimer(2 * time.Second)
            <- t.C
            e.writeHelpMenu(HELP_MESSAGE)
        }
    }
}

