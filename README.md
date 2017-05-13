# abn-amro
Community based ABN Amro bank report analyzer

**UPD**
ABN Amro released mobile application called [Grip](https://www.abnamro.nl/nl/prive/apps/grip/index.html). It is in Dutch by the idea the same as this one.

**ALPHA**
To try it:
* install [docker](https://docker.io);
* go to your ABN Amro online banking account and export your transactions in TEXT format;
* rename this file to `input.TAB` and place in the root of this project;
* do `./run`

**Output example:**

    Stats between: 2015-0-0 and 2015-0-0
    -----------------------------------------------
    food: -0.27
    rent: -0.33
    g/w/e/i: 0.00
    insurance: -0.42
    abonements: -0.02
    hema/blokker/other: -0.50
    medicine: -0.98
    electronics: -0.60
    clothes: -0.58
    deposits: -0.00
    other: -0.69
    salary: 0.36
    unknown_income: 0.00
    income: 0.35
    rest: -0.03

You can alter regular expressions and configure output in that way in `app.go` file.
This will be improved someday.

MIT license
