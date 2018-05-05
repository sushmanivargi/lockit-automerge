package util

import (
  "log"
  )

// ArrayContainsString verifies if `e` is present in the `s` array
func ArrayContainsString(s []string, e string) bool {
  if len(s) == 0 {
    return false
  }
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

// ArrayContainsInt verifies if `e` is present in the `s` array
func ArrayContainsInt(s []int, e int) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func DeleteArrayElement(s []string, e string) []string{
  var idx int
  for i, a := range s {
    if a == e {
      idx = i
      break
    }
  }
  s = append(s[:idx], s[idx+1:]...)
  return s
}

func IsError(err error) bool {
  if err != nil {
   log.Printf("[ERROR]: %s", err.Error())
  }
  return (err != nil)
}
