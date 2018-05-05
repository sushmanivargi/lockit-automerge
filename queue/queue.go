package queue

import (
  "bufio"
  "fmt"
  "os"

  "github.com/coupa/lockit-automerge/util"
)

const lockitRetryQueuePath = "lockit_retry_queue.txt"

func read() ([]string, error){
  file, err := os.OpenFile(lockitRetryQueuePath, os.O_CREATE|os.O_RDONLY, 0666)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  var items []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    items = append(items, scanner.Text())
  }
  return items, scanner.Err()
}

func write(items []string) (err error){
  f, err := os.OpenFile(lockitRetryQueuePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
  if util.IsError(err) {return}
  defer f.Close()
 
  w := bufio.NewWriter(f)
  for _, item := range items {
    fmt.Fprintln(w, item)
  }
  return w.Flush()
}

func Enqueue(item string) (err error){
  // If the file doesn't exist, create it, or append to the file
  f, err := os.OpenFile(lockitRetryQueuePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  if util.IsError(err) {return}
  defer f.Close()
  _, err = fmt.Fprintln(f, item)
  return
}

func Dequeue()(string, error) {
  items, err := read()
  if util.IsError(err) {return "", err}
  newItems := append(items[:0], items[1:]...)
  err = write(newItems)
  if util.IsError(err) {return "", err}
  return items[0], err
}

func DeleteItem(item string)(err error) {
  items, err := read()
  if util.IsError(err) {return}
  items = util.DeleteArrayElement(items, item)
  err = write(items)
  return
}

func IsEmpty()(bool){
  items, _ := read()
  if len(items) == 0 {
    return true
  }
  return false
}
