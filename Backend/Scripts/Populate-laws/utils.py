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
    cmd = f'curl "{url}" -o "{output_file}"'
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