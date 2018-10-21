package main

import (
    "strings"
)

var waiters []string
var fields []string
var input string

// Add id to the end of waiters
func Add_waiter (id string) { waiters = append (waiters, id) }

// Remove element from waiters with val = id if it exists. Return true if value 
// removed was the first element
func Remove_waiter (id string) bool {
    var index int
    for index = 0; index < len(waiters); index++ {
        if strings.Compare(id, waiters[index]) == 0 {
            waiters = append(waiters[:index], waiters[index+1:]...)
            break
        }
    }

    return index == 0
}

// Return first element in waiters or "" if waiters is empty
func Get_Front () string {
    if len(waiters) > 0 { return waiters[0] }
    return ""
}
