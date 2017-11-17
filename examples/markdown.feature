Feature: Markdown
  Scenario: Check an empty markdown file
    Given a file named "foo.md" with ""
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""

  Scenario: Check a markdown file
    Given a file named "foo.md" with:
    """
    # Title

    ## Section

    - List 1
        - Item 1
        - Item 2
    - List 2

    ```
    #!/bin/sh

    echo Hello, world!
    ```
    """
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""

  Scenario: Check a markdown file which contains a live link
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""

  Scenario: Check a markdown file which contains a dead link
    Given a file named "foo.md" with:
    """
    [The answer](https://some-say-the-answer-is-42.com)
    """
    When I run `linkcheck foo.md`
    Then the exit status should be 1
    And the stderr should contain "ERROR"

  Scenario: Check a markdown file which contains a live link in verbose mode
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `linkcheck -v foo.md`
    Then the stderr should contain "OK"

  Scenario: Check a markdown file which contains a live link in verbose mode with a long option
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `linkcheck --verbose foo.md`
    Then the stderr should contain "OK"
