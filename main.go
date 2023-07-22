package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"html/template"
)

// cool tool for parsing json: https://mholt.github.io/json-to-go/
// you can uncomment if you want to take those fields (remember to add them in the loops and the html.layout file)
type NYTbooksresponse struct {
	Results      struct {
	//	ListName                 string `json:"list_name"`
		ListNameEncoded          string `json:"list_name_encoded"`
	//	BestsellersDate          string `json:"bestsellers_date"`
		PublishedDate            string `json:"published_date"`
			Books                    []struct {
			Rank               int    `json:"rank"`
			RankLastWeek       int    `json:"rank_last_week"`
		//	PrimaryIsbn10      string `json:"primary_isbn10"`
		//	PrimaryIsbn13      string `json:"primary_isbn13"`
			Publisher          string `json:"publisher"`
		//	Description        string `json:"description"`
			Title              string `json:"title"`
			Author             string `json:"author"`
			AmazonProductURL   string `json:"amazon_product_url"`
		//	AgeGroup           string `json:"age_group"`
			BookReviewLink     string `json:"book_review_link"`
			FirstChapterLink   string `json:"first_chapter_link"`
		//	SundayReviewLink   string `json:"sunday_review_link"`
			ArticleChapterLink string `json:"article_chapter_link"`
		//	Isbns              []struct {
		//		Isbn10 string `json:"isbn10"`
		//		Isbn13 string `json:"isbn13"`
		//	} `json:"isbns"`
			BuyLinks []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"buy_links"`
			BookURI string `json:"book_uri"`
		} `json:"books"`
	} `json:"results"`
}

func main(){
	//exampleURL := "https://api.nytimes.com/svc/books/v3/lists/current/combined-print-and-e-book-fiction.json?api-key={your_key_goes_here}"

	// Builds URL from base and resource plus API key
	baseURL := "https://api.nytimes.com"
	resource := "/svc/books/v3/lists/current/combined-print-and-e-book-fiction.json"
	params := url.Values{}
	params.Add("api-key", "{your_key_goes_here}")
	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	// make GET call from the above
	response, err := http.Get(urlStr)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := NYTbooksresponse{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}

//	fmt.Printf("%+v\n", data)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "ListNameEncoded\tPublishedDate\t")

	for _, book := range data.Results.Books {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t",
			data.Results.ListNameEncoded, data.Results.PublishedDate))

		fmt.Fprintln(w, fmt.Sprintf("Book: Rank\tRankLastWeek\tPublisher\tTitle\tAuthor\tAmazonProductURL\tBookReviewLink\tFirstChapterLink\tArticleChapterLink\tBookURI\t"))
		fmt.Fprintln(w, fmt.Sprintf("     %d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t",
			book.Rank,
			book.RankLastWeek, 
			book.Publisher, 
			book.Title, 
			book.Author, 
			book.AmazonProductURL, 
			book.BookReviewLink, 
			book.FirstChapterLink, 
			book.ArticleChapterLink,
			book.BookURI,
				))

		// In case the ISBN interests you
//		for _, isbn := range book.Isbns {
//			fmt.Fprintln(w, fmt.Sprintf("Isbn:
//				 Isbn10\tIsbn13\t"))
//			fmt.Fprintln(w, fmt.Sprintf("     %s\t%s\t\n", isbn.Isbn10, isbn.Isbn13))
//		}
		for _, link := range book.BuyLinks {
			fmt.Fprintln(w, fmt.Sprintf("BuyLinks: Name\tURL\t"))
			fmt.Fprintln(w, fmt.Sprintf("         %s\t%s\t\n", link.Name, link.URL))
		}
	}

	fmt.Fprintln(w)
	w.Flush()


	// Preparing the template
	tmpl := template.Must(template.ParseFiles("layout.html"))
    
	// Creating the output file
	out, err := os.Create("output.html")
	if err != nil {
		log.Println("Cannot create file", err)
		return
	}
	defer out.Close()

	// Execute the template, writing to out, with data.
	err = tmpl.Execute(out, data)
	if err != nil {
		log.Println("Execute: ", err)
	        return
	}
}
