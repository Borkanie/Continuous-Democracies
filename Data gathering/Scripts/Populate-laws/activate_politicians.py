from bs4 import BeautifulSoup
from utils import get, createtempDir, removeTempDir, log, move_first_word_to_end 
from dataBase_interaction import replace_active_politicians_by_names, set_politician_active_by_name


def activate_all_groups(idl: int = 1, start_idg: int = 1) -> int:
    """
    Iterate group pages and set politicians found there to Active in the DB.

    Stops when the page for an `idg` does not contain the marker div.
    Returns total number of DB rows updated.
    """
    createtempDir()
    idg = start_idg
    total_updated = 0
    try:
        all_names = []
        while True:
            url = f"https://www.cdep.ro/pls/parlam/structura2015.gp?idl={idl}&idg={idg}"
            file_name = f"group_{idg}.html"
            html = get(url=url, file=file_name)
            soup = BeautifulSoup(html, "html.parser")
            grp = soup.select_one("div.grup-parlamentar-list.grupuri-parlamentare-list")
            if not grp:
                log(f"No group div found for idg={idg}; stopping iteration.")
                break

            # collect names from the group anchors
            for a in grp.select("a"):
                text = a.get_text(strip=True)
                if text:
                    all_names.append(move_first_word_to_end(text))

            idg += 1

        # deduplicate across all groups
        unique_names = list(dict.fromkeys(all_names))
        log(f"Activating {len(unique_names)} politicians across all groups")
        result = replace_active_politicians_by_names(unique_names)
        total_updated = result.get('activated', 0)

    except Exception as e:
        log(f"Error while activating groups: {e}")
    finally:
        removeTempDir()

    log(f"Finished group activation. Total DB rows updated: {total_updated}")
    return total_updated


if __name__ == "__main__":
    # run once from CLI
    activated = activate_all_groups()
    print(f"Total DB rows updated: {activated}")
