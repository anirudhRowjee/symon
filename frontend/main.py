
import requests
from time import sleep
from rich.console import Console
from rich.theme import Theme
from rich.progress import track

while True:
    data = requests.get("http://localhost:1337/metrics")
    out = data.json()
    # Rich theme to colorize in the terminal
    cve_theme = Theme
    (
    {
        "critical": "red",
        "high": "blue",
        "medium": "yellow",
        "low": "green",
        "unknown": "white",
    }
    )
    #Progress Bar 

    for step in track(range(100), description="Progress Bar "):
        sleep(0.2)

    

    console = Console()
    with console.status("[bold green]Fetching data...\n") as status:
        while out:
            #num = out.pop(0)
            sleep(1)
            console.log(f"[green]Finish fetching data[/green] {out}\n")

        console.log(f'[bold][red]Done!')

    print(out)



    
