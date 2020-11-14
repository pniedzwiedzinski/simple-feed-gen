package main

import (
  "fmt"
  "github.com/gorilla/feeds"
  "github.com/urfave/cli/v2"
  "log"
  "os"
  "time"
  "path/filepath"
  "strings"
  "io/ioutil"
)

func bootstrapFeed(title string, link string, desc string, author string, email string) *feeds.Feed {
  now := time.Now()
  feed := &feeds.Feed{
    Title: title,
    Link: &feeds.Link{Href: link},
    Description: desc,
    Author: &feeds.Author{Name: author, Email: email},
    Created: now,
  }
  return feed
}

func findGmiFiles(path string) []string {
  var files []string

  _ = filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
    if err != nil {
      log.Fatal(err)
    }
    if filepath.Ext(fpath) == ".gmi" && filepath.Base(fpath) != "index.gmi" {
      files = append(files, fpath)
    }
    return nil
  })
  return files
}

/* Embed each line in `p` html tag, so feed readers can easily display it
*/
func paragraphify(fileContent string) string {
  paragraphs := strings.Split(fileContent, "\n")
  return fmt.Sprintf("<p>%s</p>", strings.Join(paragraphs, "</p><p>"))
}


func getFileDate(fpath string) time.Time {
  info, err := os.Stat(fpath)
  if err != nil {
    log.Fatal(err)
  }
  return info.ModTime()
}

func feedEntry(fpath string, folderRoot string, url string) (*feeds.Item, error) {
  title := strings.TrimSuffix(filepath.Base(fpath), filepath.Ext(fpath))
  relativePath := strings.TrimPrefix(fpath, folderRoot)
  link := fmt.Sprintf("%s%s", url, relativePath)
  data, err := ioutil.ReadFile(fpath)
  if err != nil {
    return nil, err
  }
  desc := paragraphify(string(data))
  item := &feeds.Item {
    Title: title,
    Link: &feeds.Link{Href: link},
    Description: desc,
    Created: getFileDate(fpath),
  }
  return item, nil
}


func main() {
  app := &cli.App{
    Name:  "simple-feed-gen",
    Usage: "[Gemini path to blog] [folder with .gmi blog entries]",
    Flags: []cli.Flag {
      &cli.StringFlag {
        Name: "title",
        Aliases: []string{"t"},
        Value: "feed",
        Usage: "feed title",
      },
      &cli.StringFlag {
        Name: "description",
        Aliases: []string{"d"},
        Value: "",
        Usage: "description of your feed",
      },
      &cli.StringFlag {
        Name: "author",
        Aliases: []string{"a"},
        Value: "",
        Usage: "author of the feed",
      },
      &cli.StringFlag {
        Name: "email",
        Aliases: []string{"e"},
        Value: "",
        Usage: "author email",
      },
    },
    Action: func(c *cli.Context) error {
      url := c.Args().Get(0)
      path := c.Args().Get(1)
      if stat, err := os.Stat(path); err == nil && stat.IsDir() {
        feed := bootstrapFeed(
          c.String("title"),
          url,
          c.String("description"),
          c.String("author"),
          c.String("email"),
        )
        files := findGmiFiles(path)
        if len(files) == 0 {
          log.Fatal(fmt.Sprintf("Couldn't find .gmi files in %s", path))
        }
        for _, f := range files {
          entry, err := feedEntry(f, path, url)
          if err != nil {
            return err
          }
          feed.Add(entry)
        }
        atom, err := feed.ToAtom()
        if err != nil {
          return err
        }
        fmt.Println(atom)
        return nil
      }
      log.Fatal("Pass valid path")
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
