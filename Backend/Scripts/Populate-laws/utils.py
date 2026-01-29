from time import timezone
from log_writer import file_log
import os
import subprocess
from charset_normalizer import from_path
from bs4 import BeautifulSoup
from shutil import rmtree
from datetime import datetime, timezone, timedelta
import locale
import re

def normalize_romanian_reversed_name(full_name: str) -> str:
    if not full_name:
        return ""

    words = full_name.strip().split()

    # Last word is always the FAMILY NAME and always full caps
    family_raw = words[-1]
    first_names_raw = words[:-1]

    # Safety check: ensure last word is really the FAMILY name
    if not family_raw.isupper():
        # If not full caps, fallback: treat whole name normally
        # Capitalize all words and return unchanged order
        return " ".join(w.capitalize() for w in words)

    # Helper: Romanian-aware capitalization
    def cap(word):
        if not word:
            return word
        return word[0].upper() + word[1:].lower()

    family_name = cap(family_raw)

    first_names = " ".join(cap(w) for w in first_names_raw)

    # Return FAMILY first
    return f"{family_name} {first_names}".strip()

def generate_party_acronym(name: str) -> str:
    if not name or not isinstance(name, str):
        return ""

    words = re.split(r"\s+", name.strip())
    acronym_letters = []

    for word in words:
        clean_word = re.sub(r"[^a-zA-ZăîâșțĂÎÂȘȚ]", "", word)  # remove punctuation

        # Skip words that are ALL CAPS
        if clean_word.isupper():
            continue

        # Skip short words
        if len(clean_word) <= 3:
            continue

        acronym_letters.append(clean_word[0])

    # Uppercase final acronym
    return "".join(acronym_letters).upper()


def move_first_word_to_end(full_name: str) -> str:
    """
    Move the first word of a full name to the end.

    Examples:
        'Popa Ştefan-Ovidiu' -> 'Ştefan-Ovidiu Popa'
        'Adomnicăi Mirela Elena' -> 'Mirela Elena Adomnicăi'

    This function preserves Unicode characters (Romanian/Hungarian diacritics)
    and collapses multiple whitespace characters. If the input contains only
    one word, it is returned unchanged.
    """
    if not full_name or not isinstance(full_name, str):
        return full_name

    parts = re.split(r"\s+", full_name.strip())
    if len(parts) <= 1:
        return full_name.strip()

    reordered = parts[1:] + [parts[0]]
    return " ".join(reordered)

def deputy_exists(dom_str: str) -> str | None:
    """
    Returns True if the DOM contains a real deputy profile,
    False if the page is an empty placeholder.
    """
    soup = BeautifulSoup(dom_str, "html.parser")

    # 1) Best source: the main headline
    tag = soup.select_one("td.headline")
    if tag:
        text = tag.get_text(" ", strip=True)
        # Headline always contains: "<NAME> Sinteza ..."
        name = text.split(" Sinteza")[0].strip()
        return name

    # 2) Fallback: small profile block (left column)
    block = soup.select_one("tr[id^=itm974] td")
    if block:
        text = block.get_text(" ", strip=True)
        # Usually: "Name n. 17 apr. 1985"
        name = text.split(" n.")[0].strip()
        return name

    return None


def extract_datetime_from_dom(dom_str):
    # Set Romanian locale (fallback for Windows)
    try:
        locale.setlocale(locale.LC_TIME, "ro_RO.UTF-8")
    except:
        try:
            locale.setlocale(locale.LC_TIME, "ro_RO")
        except:
            pass  # If it fails, month names still work for most systems

    soup = BeautifulSoup(dom_str, "html.parser")
    raw_text = soup.select_one(".boxTitle h1").get_text(" ", strip=True)

    # Normalize multiple spaces
    text = re.sub(r"\s+", " ", raw_text)

    # Remove weekday
    text = text.split(",", 1)[1].strip()  
    # Remove the 'ora' keyword
    text = text.replace("ora", "").strip()  
    # Remove comma between date and time
    text = text.replace(",", "").strip()

    # Now input looks like: "13 iunie 2023 12:16"
    dt = datetime.strptime(text, "%d %B %Y %H:%M")

    # Romanian timezone (summer = UTC+3, winter = UTC+2)
    RO_STANDARD = timezone(timedelta(hours=2))
    RO_SUMMER   = timezone(timedelta(hours=3))

    # Rough DST detection
    if dt.month < 3 or dt.month > 10:
        local_tz = RO_STANDARD
    elif 4 <= dt.month <= 9:
        local_tz = RO_SUMMER
    else:
        # March or October — approximate rule (close enough for your dates)
        local_tz = RO_SUMMER if dt.month == 3 else RO_STANDARD

    dt_local = dt.replace(tzinfo=local_tz)
    dt_utc = dt_local.astimezone(timezone.utc)

    # Your required output format:
    return dt_utc.strftime("%Y-%m-%d %H:%M:%S.%f+00")

def downloadFileFormUrl(url, output_file):
    """
    Download a file from a given URL and save it to the specified output file.
    Uses curl command to perform the download. Errors are captured and logged,
    and a RuntimeError is raised if the curl command fails.
    Args:
        url (str): The URL of the file to download.
        output_file (str): The file path where the downloaded file will be saved.
    Raises:
        RuntimeError: If the curl command fails (non-zero exit code).
    Returns:
        None
    """
    cmd = f'curl -s "{url}" -o "{output_file}"'
    # internal download messages go to file
    file_log("Running curl for:", url, "->", output_file)
    #os.system(cmd)
    result = subprocess.run(
        cmd,
        shell=True,
        stdout=subprocess.DEVNULL,   # hide normal output
        stderr=subprocess.PIPE,      # capture errors
        text=True
    )

    if result.stderr:
        print(result.stderr)

    if result.returncode != 0:
        raise RuntimeError(f"curl command failed with exit code {result.returncode} for URL: {url}")

def get(url, file= "temp_response.txt", keep = False) -> str:
        """
        Downloads the content from the specified URL and saves it to a temporary file, then reads and returns the file's content as a string and deletes the file at the end
        Args:
            url (str): The URL to download content from.
            file (str, optional): The name of the temporary file to save the response. Defaults to "temp_response.txt".
        Returns:
            str: The content retrieved from the URL.
        Raises:
            Exception: If there is an error removing the temporary file, logs the failure but does not raise.
        """
        temp_dir = createtempDir()
        output_file = os.path.join(temp_dir, file)
        if not os.path.exists(output_file):
            downloadFileFormUrl(url, output_file)

        # Read the response from the file
        with open(output_file, 'r', encoding=from_path(output_file).best().encoding, errors='replace') as f:
            response_text = f.read()
        
        # Delete the temporary file
        try:
            if  not keep:
                os.remove(output_file)
        except Exception:
            file_log("Failed to remove temp file:", output_file)
        return response_text

def createtempDir():
    temp_dir = os.path.join(os.getcwd(), "temp")
    os.makedirs(temp_dir, exist_ok=True)
    return temp_dir

def removeTempDir():
    try:
        rmtree(os.path.join(os.getcwd(), "temp"))
    except Exception:
        print("Failed to remove temp directory")
        pass

def first_if_list(value):
    """
    Returns the first element of a list or tuple if the input is a list or tuple; otherwise, returns the input value itself.
    If the list or tuple is empty, returns None.

    Args:
        value: The input value, which may be a list, tuple, or any other type.

    Returns:
        The first element of the list or tuple, None if empty, or the value itself if not a list or tuple.
    """
    if isinstance(value, (list, tuple)):
        return value[0] if value else None
    return value