package main

import (
    "strings"
    "strconv"
    "time"
    "bufio"
    "regexp"
    "fmt"
    "os"
)

var rules = map[string]map[string]*regexp.Regexp{
    "food": map[string]*regexp.Regexp{
        "ALBERT HEIJN": regexp.MustCompile(`ALBERT HEIJN`),
        "Jumbo": regexp.MustCompile(`Jumbo[^a-zA-Z0-9-]`),
    },
    "hema/blokker/other": map[string]*regexp.Regexp{
        "HEMA": regexp.MustCompile(`HEMA`),
        "BEDWORLD": regexp.MustCompile(`BEDWORLD`),
    },
    "electronics": map[string]*regexp.Regexp{
        "Apple": regexp.MustCompile(`APPLE STORE`),
        "Coolblue": regexp.MustCompile(`Coolblue`),
        "Wiggle": regexp.MustCompile(`Wiggle`),
    },
    "clothes": map[string]*regexp.Regexp{
        "Decathlon": regexp.MustCompile(`Decathlon`),
        "FRONT RUNNER": regexp.MustCompile(`FRONT RUNNER`),
    },
    "medicine": map[string]*regexp.Regexp{
        "APOTHEEK": regexp.MustCompile(`APOTHEEK`),
    },
    "deposits": map[string]*regexp.Regexp{
        "Wallet": regexp.MustCompile(`IBAN\/NL14ABNA0620233052`),
    },
    "rent": map[string]*regexp.Regexp{
        "Dijwater 225": regexp.MustCompile(`IBAN\/NL18ABNA0459968513`),
    },
    "g/w/e/i": map[string]*regexp.Regexp{
        "Vodafone": regexp.MustCompile(`CSID\/NL39ZZZ302317620000`),
        "Energiedirect": regexp.MustCompile(`CSID\/NL71ABNA0629639183`),
    },
    "abonements": map[string]*regexp.Regexp{
        "Sport natural": regexp.MustCompile(`SPORT NATURAL`), // IBAN\/NL07ABNA0584314078
    },
    "insurance": map[string]*regexp.Regexp{
        "Dijwater 225": regexp.MustCompile(`CSID\/NL94MNZ505448100000`),
    },
    "salary": map[string]*regexp.Regexp{
        "Textkernel": regexp.MustCompile(`IBAN\/NL27INGB0673657841`),
    },
}

var displayOrder = []string{
    "food",
    "rent",
    "g/w/e/i",
    "insurance",
    "abonements",
    "hema/blokker/other",
    "medicine",
    "electronics",
    "clothes",
    "deposits",
    "other",
    "salary",
    "unknown_income",
    "income",
    "rest",
}

// var reClarify = regexp.MustCompile(`^\/[^\/]+/[^\/]+/[^\/]+/[^\/]+/[^\/]+/[^\/]+/[^\/]+/([^\/]+)/[^\/]+/[^\/]+/[^\/]+/[^\/]+$`)
// var reClarify2 = regexp.MustCompile(`^[^ ]+\s+[^ ]+\s+\d\d\.\d\d\.\d\d\/\d\d\.\d\d\s+(.*),[A-Z0-9 ]+$`)
// var reClarify3 = regexp.MustCompile(`^[^ ]+\s+\d\d-\d\d-\d\d \d\d:\d\d\s+[^ ]+\s+(.*),[A-Z0-9 ]+$`)

type Payment struct {
    accountId int64
    currency string
    date time.Time
    balanceBefore float64
    balanceAfter float64
    anotherDate time.Time
    amount float64
    description string
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


// func clarifyName(name string) (string, error) {
//     matches := reClarify.FindStringSubmatch(name)
//     if len(matches) > 1 {
//         return matches[1], nil
//     }
//
//     matches = reClarify2.FindStringSubmatch(name)
//     if len(matches) > 1 {
//         return matches[1], nil
//     }
//
//     matches = reClarify3.FindStringSubmatch(name)
//     if len(matches) > 1 {
//         return matches[1], nil
//     }
//
//     return ``, fmt.Errorf(`Not parsable: %s`, name)
// }

func main() {
    file, err := os.Open("input.TAB")
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    re := regexp.MustCompile(`^(\d+)\t([A-Za-z]+)\t(\d+)\t(\d+,\d+)\t(\d+,\d+)\t(\d+)\t([-]?\d+,\d+)\t(.*?)$`);
    payments := make([]*Payment, 0, 5)
    for scanner.Scan() {
        txt := scanner.Text()
        matched := re.FindAllStringSubmatch(txt, 1)

        p := &Payment{};

        p.accountId, err = strconv.ParseInt(matched[0][1], 10, 64)
        check(err)

        p.currency = matched[0][2]

        p.date, err = time.Parse("20060102", matched[0][3])
        check(err)

        p.balanceBefore, err = strconv.ParseFloat(strings.Replace(matched[0][4], ",", ".", 1), 64)
        check(err)

        p.balanceAfter, err = strconv.ParseFloat(strings.Replace(matched[0][5], ",", ".", 1), 64)
        check(err)

        p.anotherDate, err = time.Parse("20060102", matched[0][6])
        check(err)

        p.amount, err = strconv.ParseFloat(strings.Replace(matched[0][7], ",", ".", 1), 64)
        check(err)

        p.description = matched[0][8]

        payments = append(payments, p)
    }

    check(scanner.Err())

    results := map[string]float64{
        `supermarkets`: 0,
        `other`: 0,
        `income`: 0,
        `unknown_income`: 0,
        `rest`: 0,
    }

    var firstPaymentDate time.Time
    var lastPaymentDate time.Time
    i := 0
    paymentsLen := len(payments)
    for _, p := range payments {
        switch i {
            case 0: firstPaymentDate = p.date
            case paymentsLen - 1: lastPaymentDate = p.date
        }
        other := true
        for rulesName, rulesList := range rules {
            for _, re := range rulesList {
                if re.MatchString(p.description) {
                    if _, ok := results[rulesName]; ok {
                        results[rulesName] += p.amount
                    } else {
                        results[rulesName] = p.amount
                    }
                    other = false
                    break
                }
            }
        }

        if other && p.amount <= 0 {
            results[`other`] += p.amount
        }

        if other && p.amount > 0 {
            results[`unknown_income`] += p.amount
        }

        if p.amount > 0 {
            results[`income`] += p.amount
        }

        results[`rest`] += p.amount
        i += 1
    }


    fmt.Printf("Stats between: %s and %s\n", firstPaymentDate.Format(`2006-01-02`), lastPaymentDate.Format(`2006-01-02`))
    fmt.Println("-----------------------------------------------")
    for _, rulesName := range displayOrder {
        fmt.Printf("%s: %.2f\n", rulesName, results[rulesName])
    }
}