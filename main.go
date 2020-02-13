// testapp project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/alexeyco/simpletable"
)

type pkg_struct struct {
	Package       string
	Version       string
	Depends       string
	Section       string
	Architecture  string
	InstalledSize string
	Filename      string
	Size          string
	SHA256sum     string
	Description   string
	Feed          string
}

type pkg_list_struct struct {
	pkgs_list []pkg_struct
	Name      string
}

type branches struct {
	Name     string
	url      string
	feeds    []string
	pkg_list []pkg_struct
}

/*type pkg_db struct {
}*/

func download_page(url string) (page string) {
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// show the HTML code as a string %s
	page = string(html)

	return page
}

func update_pkg_struct(name, value string, pkg *pkg_struct) {
	//pkg_s:=&pkg

	switch name {
	case "Package":
		pkg.Package = value
	case "Version":
		pkg.Version = value
	case "Depends":
		pkg.Depends = value
	case "Section":
		pkg.Section = value
	case "Architecture":
		pkg.Architecture = value
	case "InstalledSize":
		pkg.InstalledSize = value
	case "Filename":
		pkg.Filename = value
	case "Size":
		pkg.Size = value
	case "SHA256sum":
		pkg.SHA256sum = value
	case "Description":
		pkg.Description = value
	}
}

// Parse package to structure
func parse_pkg(text string) (pkg pkg_struct) {
	var tmp_pkg pkg_struct
	for _, el := range strings.Split(text, "\n") {
		val := strings.Split(el, ":")
		if len(val) == 2 {
			name := strings.TrimSpace(val[0])
			value := strings.TrimSpace(val[1])
			//Update package structures
			update_pkg_struct(name, value, &tmp_pkg)
		}
	}
	return tmp_pkg
}

// Parse list of packages from Packages file to array
func parse_pkg_list(text string) (pkg_list []pkg_struct) {
	for _, el := range strings.Split(text, "\n\n") {
		tmp_pkg := parse_pkg(el)
		pkg_list = append(pkg_list, tmp_pkg)
	}
	return pkg_list
}

func init_branch_list(url_temp string, branch_list []branches) {
	for index, item := range branch_list {
		branch_list[index].url = fmt.Sprintf(url_temp, item.Name, "%s")
	}
}

func print_pkg_tbl(pkg []pkg_struct) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Version"},
			{Align: simpletable.AlignCenter, Text: "Section"},
		},
	}
	for _, row := range pkg {
		r := []*simpletable.Cell{
			{Text: fmt.Sprintf("%s", row.Package)},
			{Text: row.Version},
			{Text: row.Section},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleRounded)
	fmt.Println(table.String())
}

func find_pkg_comp(pkg_name string, pkg_lists []pkg_list_struct) {
	fmt.Printf("testik")
}

func find_pkg(pkg_name string, pkg_list []pkg_struct) (pkg_list_ret []pkg_struct) {
	//var ret_pkg_list []pkg_struct
	for _, item := range pkg_list {
		if strings.Contains(strings.ToLower(item.Package), pkg_name) {
			pkg_list_ret = append(pkg_list_ret, item)
		}
	}
	return pkg_list_ret
}

func find_dep(pkg_name string, pkg_list []pkg_struct) (pkg_list_ret []pkg_struct) {
	//var ret_pkg_list []pkg_struct
	for _, item := range pkg_list {
		if strings.Contains(strings.ToLower(item.Depends), pkg_name) {
			pkg_list_ret = append(pkg_list_ret, item)
		}
	}
	return pkg_list_ret
}

/*
func downloadFile(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func downloadMultipleFiles(urls []string) ([][]byte, error) {
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := downloadFile(URL)
			if err != nil {
				errch <- err
				done <- nil
				return
			}
			done <- b
			errch <- nil
		}(URL)
	}
	bytesArray := make([][]byte, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {
		bytesArray = append(bytesArray, <-done)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}
	return bytesArray, err
}
*/
func download_all(branch_list []branches) {
	for index, item := range branch_list {
		fmt.Printf("%s\n", item.Name)
		for _, url_tmp := range item.feeds {

			full_url := fmt.Sprintf(item.url, url_tmp)
			aa := download_page(full_url)
			bb := parse_pkg_list(aa)
			branch_list[index].pkg_list = append(branch_list[index].pkg_list, bb...)
		}
	}
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func download_branch_list(branch_list []branches, branch_list_name []string) {
	for index, item := range branch_list {
		fmt.Printf("%s\n", item.Name)
		_, found := Find(branch_list_name, item.Name)
		if found {
			for _, url_tmp := range item.feeds {
				full_url := fmt.Sprintf(item.url, url_tmp)
				aa := download_page(full_url)
				bb := parse_pkg_list(aa)
				branch_list[index].pkg_list = append(branch_list[index].pkg_list, bb...)
			}
		}

	}
}

func main() {

	var branch_name string
	var branch_index int

	branch_list := []branches{
		branches{Name: "hbd", url: "https://repo.turris.cz/hbd/omnia/packages/%s/Packages", feeds: []string{"turrispackages", "packages", "base"}},
		branches{Name: "hbl", url: "https://repo.turris.cz/hbl/omnia/packages/%s/Packages", feeds: []string{"turrispackages", "packages", "base"}},
		branches{Name: "hbk", url: "https://repo.turris.cz/hbk/omnia/packages/%s/Packages", feeds: []string{"turrispackages", "packages", "base"}},
		branches{Name: "hbs", url: "https://repo.turris.cz/hbs/omnia/packages/%s/Packages", feeds: []string{"turrispackages", "packages", "base"}},
		branches{Name: "hbt", url: "https://repo.turris.cz/hbt/omnia/packages/%s/Packages", feeds: []string{"turrispackages", "packages", "base"}},
		branches{Name: "3x", url: "https://repo.turris.cz/omnia/packages/%s/Packages", feeds: []string{"turrispackages", "packages", "base"}},
	}

	parser := argparse.NewParser("print", "Prints provided string to stdout")

	// Create string flag
	f := parser.String("f", "find", &argparse.Options{Required: false, Help: "Find package"})
	d := parser.String("d", "dep", &argparse.Options{Required: false, Help: "Find package dependency"})

	hbl := parser.Flag("l", "hbl", &argparse.Options{Help: "Search in HBL"})
	hbd := parser.Flag("d", "hbd", &argparse.Options{Help: "Search in HBD"})
	hbk := parser.Flag("k", "hbk", &argparse.Options{Help: "Search in HBK"})
	hbt := parser.Flag("t", "hbt", &argparse.Options{Help: "Search in HBT"})
	hbs := parser.Flag("s", "hbs", &argparse.Options{Help: "Search in HBS"})
	old := parser.Flag("3", "3x", &argparse.Options{Help: "Search in 3x"})
	comp := parser.Flag("c", "compare", &argparse.Options{Help: "Compare versions"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	var branches []string

	if *hbl {
		branch_name = "hbl"
		branch_index = 1
		branches = append(branches, branch_name)
	} else if *hbd {
		branch_index = 0
		branch_name = "hbd"
		branches = append(branches, branch_name)
	} else if *hbk {
		branch_index = 2
		branch_name = "hbk"
		branches = append(branches, branch_name)
	} else if *hbt {
		branch_name = "hbt"
		branch_index = 4
		branches = append(branches, branch_name)
	} else if *hbs {
		branch_name = "hbs"
		branch_index = 3
		branches = append(branches, branch_name)
	} else if *old {
		branch_name = "3x"
		branch_index = 5
		branches = append(branches, branch_name)
	} else {
		branch_name = "hbd"
		branch_index = 0
		branches = append(branches, branch_name)
	}

	if *comp {
		fmt.Println("Compare packages")
	}

	// Finally print the collected string
	if len(*f) > 0 {
		download_branch_list(branch_list, []string{branch_name})
		aa := branch_list[branch_index].pkg_list
		cc := find_pkg(*f, aa)
		print_pkg_tbl(cc)
	}
	if len(*d) > 0 {
		txt := download_page("https://repo.turris.cz/hbd/omnia/packages/packages/Packages")
		aa := parse_pkg_list(txt)
		cc := find_dep(*d, aa)
		print_pkg_tbl(cc)
	}
	fmt.Println(*f)

}
