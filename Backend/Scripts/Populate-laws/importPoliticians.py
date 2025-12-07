from utils import createtempDir, removeTempDir
from dataBaseInteraction import getLaws, getPoliticans, insertPolitician, insertNewParty, getParties
from uuid import uuid4
from utils import file_log, get, deputy_exists, generate_party_acronym
from bs4 import BeautifulSoup

parties = getParties()
politicians = getPoliticans()

class Politician:
    def __init__(self, imageUrl:str, name: str, partyId: str, gender: int, active: bool):
        self.id = str(uuid4())
        self.gender = gender
        self.imageUrl = imageUrl
        self.name = name
        self.active = active
        self.partyId = partyId

def checkIfPartyExists(partyName: str) -> str | None:
    matching_parties = [party for party in parties if partyName in party["Name"]]
    if len(matching_parties) == 0:
        return None
    else:
        return matching_parties[0]["Id"]

def normalize_name(name: str) -> str:
    """
    In the input, the FAMILY name is ALWAYS last and ALL-CAPS.
    We transform:
        'Alexandru ALBU' → 'Alexandru Albu'
        'Andrei-Ionuţ TESLARIU' → 'Andrei-Ionuţ Teslariu'
    """
    parts = name.split()
    if not parts:
        return name

    last = parts[-1]
    if last.isupper():
        last = last.capitalize()
        return " ".join(parts[:-1] + [last])

    return name

def extract_politician(dom_str: str, year: int) -> Politician | None:
    soup = BeautifulSoup(dom_str, "html.parser")

    # ------------------------------------------------------
    # 1) NAME (legacy layout)
    # ------------------------------------------------------
    # First attempt: "headline" block (legacy)
    name_tag = soup.select_one("td.headline")
    if name_tag:
        raw_name = name_tag.get_text(strip=True).split("\n")[0]
    else:
        # fallback: image alt text (legacy layout)
        img_tag = soup.select_one("table img[alt]")
        raw_name = img_tag["alt"].strip() if img_tag else ""

    if not raw_name:
        return None

    name = normalize_name(raw_name.split("Sinteza")[0].strip())

    # ------------------------------------------------------
    # 2) IMAGE URL (legacy)
    # ------------------------------------------------------
        # ---------- IMAGE ----------
    # Modern: inside .profile-pic-dep
    img_tag = soup.select_one(".profile-pic-dep img")
    if not img_tag:
        # Legacy: td.menuoff with image width 150
        img_tag = soup.select_one('td.menuoff img[width="150"]')

    image_url = "https://www.cdep.ro" + img_tag["src"] if img_tag and img_tag.has_attr("src") else None
    # ------------------------------------------------------
    # 3) PARTY NAME + ACRONYM
    # ------------------------------------------------------
    # Legacy layout:
    #   <td><a href="structura.fp?...">FDSN</a></td>
    party_tag = soup.select_one('table a[href*="structura.fp"]')
    if party_tag:
        party_acronym = party_tag.get_text(strip=True)
    else:
        party_acronym = None

    # ------------------------------------------------------
    # 4) Detect full party name (legacy)
    # ------------------------------------------------------
    # Structure:
    #   <td>FDSN</td><td>-</td><td>Frontul Democrat al ...</td>
    full_party_name = None
    if party_tag and party_tag.parent:
        cells = list(party_tag.parent.parent.find_all("td"))
        if len(cells) >= 3:
            full_party_name = cells[2].get_text(strip=True)

    # If no full name found → use acronym
    partyId = checkIfPartyExists(full_party_name)
    if partyId is None:
        file_log(f"Party not found for politician {name}, party acronym: {full_party_name}", alsoPrint=True)
        party_acronym = party_acronym if party_acronym is not None else generate_party_acronym(full_party_name)
        partyId = insertNewParty(party_acronym, full_party_name)
        if partyId is None:
            file_log(f"Failed to insert new party for politician {name}, party acronym: {full_party_name}", alsoPrint=True)
            return None
        
    # ------------------------------------------------------
    # 5) GENDER HEURISTIC
    # ------------------------------------------------------
    last_name = name.split(" ")[0]
    if last_name.endswith("a"):
        gender = 1
    else:
        gender = 2

    # ------------------------------------------------------
    # 6) ACTIVE FLAG
    # ------------------------------------------------------
    active = True

    # ------------------------------------------------------
    # 7) CREATE Politician OBJECT
    # ------------------------------------------------------
    return Politician(
        imageUrl=image_url,
        name=name,
        partyId=partyId,
        gender=gender,
        active=active
    )


def importPoitician(politicianId: int, year: int) -> bool:
    url = f"https://www.cdep.ro/pls/parlam/structura.mp?idm={politicianId}&cam=2&leg={year}"
    text = get(url=url, file=f"politician_{politicianId}_{year}.html")
    if deputy_exists(text):
        politician = extract_politician(text, year)
        if politician is None:
            file_log(f"Failed to extract politician with ID {politicianId} for year {year}.", alsoPrint=True)
            return False
        matching = [pol for pol in politicians if pol["Name"] is politician.name]
        if len(matching) > 0:
            file_log(f"Politician {politician.name} already exists in database.", alsoPrint=True)
            return True
        else:
            insertPolitician(politician.name, politician.partyId, politician.imageUrl, politician.gender, politician.active)
    else:
        file_log(f"Politician with ID {politicianId} does not exist for year {year}.", alsoPrint=True)
        return False

def importPoliticians(startId: int, endId: int,year: int):
    """
    Imports politicians with IDs in the specified range.
    Args:
        startId (int): The starting ID of the politicians to import.
        endId (int): The ending ID of the politicians to import.
    """
    for politicianId in range(startId, endId + 1):
        try:
            print(f"Importing politician with ID: {politicianId}")
            importPoitician(politicianId, year)
        except Exception as e:
            file_log(f"Error importing politician with ID {politicianId} for year{year}: {e}", alsoPrint=True)
        

if __name__ == "__main__":
    createtempDir()
    for year in range(1990, 2025, 4):
        importPoliticians(1, 333, year)
    removeTempDir()