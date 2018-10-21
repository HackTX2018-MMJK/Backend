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
    index := Get_position(id)
    if index != -1 { waiters = append(waiters[:index], waiters[index+1:]...) }
    return index == 0
}

// Return first element in waiters or "" if waiters is empty
func Get_front () string {
    if len(waiters) > 0 { return waiters[0] }
    return ""
}

// Return the position of the id in waiters, -1 otherwise
func Get_position (id string) int {
    var index int                                                               
    for index = 0; index < len(waiters); index++ {                              
        if strings.Compare(id, waiters[index]) == 0 { return index }
    }
    return -1
}

// Return the number of elements in waiters
func Get_length () int { return len(waiters) }
