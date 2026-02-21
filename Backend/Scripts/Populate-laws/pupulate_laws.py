import re
import requests
import os
import time
import unicodedata
import shutil
import subprocess
import uuid
from PIL import Image
import pytesseract
import json
from typing import Optional
from log_writer import file_log
from PyPDF2 import PdfReader
from typing import Optional
from dataBase_interaction import delete_law, getUnpopulatedLaws, update_law, update_law_description, promulgareKey, adoptedKey, motivatieKey, ougKey, VoteIdKey, DescriptionKey, TitleKey, IdKey
from utils import downloadFileFormUrl, get

def extract_text_from_selectable_pdf(pdf_path: str) -> str:
    """
    Extract text from a PDF with selectable text (no OCR).

    Args:
        pdf_path: Path to the PDF file.

    Returns:
        Text content of the PDF as a single string.
    """
    try:
        reader = PdfReader(pdf_path)
        text = []
        file_log(f"Opened PDF {pdf_path} with {len(reader.pages)} pages.")
        for page in reader.pages:
            file_log(f"Extracting text from page {reader.pages.index(page) + 1}")
            page_text = page.extract_text()
            if page_text:
                text.append(page_text)
        file_log(f"Extracted text from PDF {pdf_path}, total length: {sum(len(t) for t in text)} characters.")
        return "\n".join(text)
    except Exception as e:
        print(f"Error reading PDF {pdf_path}: {e}")
        return ""


def ask_chatgpt_json(prompt: str, api_key: Optional[str] = None, model: str = "gpt-4.1-nano") -> Optional[dict]:
    """
    Send a prompt to the OpenAI ChatGPT API and extract JSON from the response.
    
    Args:
        prompt: The prompt to send (expects ChatGPT to respond with JSON).
        api_key: Your OpenAI API key.
        model: GPT model to use (default: gpt-4).
    
    Returns:
        Parsed JSON as a dict, or None if parsing fails.
    """
    # Prefer the provided api_key, fall back to environment variable OPENAI_API_KEY
    if not api_key:
        api_key = os.getenv("OPENAI_API_KEY")

    if not api_key:
        file_log("OpenAI API key not set. Set OPENAI_API_KEY in environment or pass api_key parameter.")
        return None

    url = "https://api.openai.com/v1/chat/completions"

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json"
    }
    
    data = {
        "model": model,
        "messages": [{"role": "user", "content": prompt}],
        "temperature": 0
    }
    
    try:
        response = requests.post(url, headers=headers, json=data, timeout=30)
        response.raise_for_status()
        resp_json = response.json()
        
        # Extract the assistant's text content
        text = resp_json["choices"][0]["message"]["content"]
        
        # Find JSON in the text (supports ```json ... ``` or plain JSON)
        json_match = re.search(r"```json\s*(\{.*?\})\s*```", text, re.DOTALL)
        if json_match:
            json_str = json_match.group(1)
        else:
            # fallback: try to parse entire text
            json_str = text.strip()
        
        return json.loads(json_str)
    
    except Exception as e:
        # Could log e if needed
        return None



def extract_text_from_pdf(pdf_path: str, poppler_path: str, tesseract_cmd: str = None, dpi: int = 300) -> str:
    """
    Extract Romanian text from a PDF using Poppler (pdfinfo + pdftoppm) and Tesseract OCR.
    
    Args:
        pdf_path: Path to the PDF file.
        poppler_path: Path to the Poppler 'bin' folder containing pdfinfo.exe and pdftoppm.exe.
        tesseract_cmd: Optional full path to tesseract executable.
        dpi: Resolution for converting PDF pages to images.
    
    Returns:
        Extracted text as a string.
    """
    temp_dir = os.path.join(os.path.dirname(pdf_path), "temp_" + uuid.uuid4().hex)
    os.makedirs(temp_dir, exist_ok=True)
    # log internal info to file
    file_log("Temporary directory created at:", temp_dir)

    try:
        pdfinfo_exe = os.path.join(poppler_path, "pdfinfo.exe")
        pdftoppm_exe = os.path.join(poppler_path, "pdftoppm.exe")

        if not os.path.exists(pdfinfo_exe) or not os.path.exists(pdftoppm_exe):
            raise FileNotFoundError("Poppler executables not found in the provided path.")

        # Step 1: Get page count
        result = subprocess.run([pdfinfo_exe, pdf_path], capture_output=True, text=True, check=True)
        page_count = 0
        for line in result.stdout.splitlines():
            if line.lower().startswith("pages:"):
                page_count = int(line.split(":", 1)[1].strip())
                break
        if page_count == 0:
            raise RuntimeError("Could not get page count from pdfinfo.")

        # Step 2: Convert PDF to images
        output_prefix = os.path.join(temp_dir, f"page_{uuid.uuid4().hex}")
        subprocess.run([
            pdftoppm_exe,
            "-r", str(dpi),
            "-png",
            pdf_path,
            output_prefix
        ], check=True, capture_output=True)

        # Collect generated images
        prefix_basename = os.path.basename(output_prefix)
        images = sorted([
            os.path.join(temp_dir, f)
            for f in os.listdir(temp_dir)
            if f.startswith(prefix_basename) and f.endswith(".png")
        ])
        if not images:
            raise RuntimeError("No images were generated from PDF.")

        # Step 3: OCR each image
        if tesseract_cmd:
            pytesseract.pytesseract.tesseract_cmd = tesseract_cmd

        text_output = []
        for img_path in images:
            text_output.append(pytesseract.image_to_string(Image.open(img_path), lang="ron"))
        file_log(f"OCR completed on {pdf_path}! extracted {len(text_output)} pages.")
        return "\n".join(text_output).strip()
    except Exception as e:
        # error prints should remain on console
        print(f"Error during PDF text extraction: {e} for file {pdf_path}")
        return None
    finally:
        # Cleanup temp files
        if os.path.exists(temp_dir):
            shutil.rmtree(temp_dir)
      

def getTextTroughOCRFromUrl(lawUrl):
    file_log("Fetching motivatie URL:", lawUrl)
    output_file = "E:\\populate-laws\\temp\\motivatie.pdf" 
    output_file = unicodedata.normalize("NFKD", output_file).replace("\\", "/")
    downloadFileFormUrl(lawUrl, output_file)
    time.sleep(1)  # allow OS to release file handle
    text = extract_text_from_pdf(output_file, "E:/poppler-25.11.0/Library/bin", "D:/tesseract-docker/tesseract.exe", dpi=300)
    try:
        os.remove(output_file)
    except Exception:
        file_log("Could not remove downloaded motivatie file:", output_file)
    return text

def fetch_and_search(url, regex_map):
    """
    Makes an HTTP GET request to the given URL, searches the response content using the provided regex map,
    and returns a dictionary of matching results.

    :param url: The URL to fetch data from.
    :param regex_map: A dictionary where keys are descriptive names and values are regex patterns.
    :return: A dictionary with the same keys as regex_map and values as lists of matches.
    """
    try:

        response = get(url)

        # Search the response content with the regex map
        matches = {}
        for key, pattern in regex_map.items():
            found = re.findall(pattern, response)
            matches[key] = found[0] if found else None
            # log search attempts internally
            file_log(f"fetch_and_search: key={key} matched_count={len(found)}")
            
        return matches

    except requests.RequestException as e:
        # error prints should remain on console
        print(f"HTTP request error: {e}")
        return {}

def getTextFromLawByForm(getlawForm):
    baseUrl = "https://www.cdep.ro"
    if getlawForm[promulgareKey]:
        lawUrl = baseUrl + "/pls/proiecte/" + getlawForm[promulgareKey]
        output_file = "E:\\populate-laws\\temp\\promulgare.pdf"
        downloadFileFormUrl(lawUrl, output_file)
        text =  extract_text_from_selectable_pdf(output_file)
        try:
            os.remove(output_file)
        except Exception:
            file_log("Could not remove downloaded promulgare file:", output_file)
        return text
    
    if getlawForm[adoptedKey]:
        lawUrl = baseUrl + getlawForm[adoptedKey]
        return getTextTroughOCRFromUrl(lawUrl)
    
    if getlawForm[ougKey]:
        lawUrl = baseUrl + getlawForm[ougKey]
        return getTextTroughOCRFromUrl(lawUrl)

    if getlawForm[motivatieKey]:
        lawUrl = baseUrl + getlawForm[motivatieKey]
        return getTextTroughOCRFromUrl(lawUrl)
    
    return None



# Example usage
def getNewNameAndDescriptionForLaw(first_law, onlyDesc: bool = False) -> Optional[dict]:
    vote_id = first_law[VoteIdKey]  # Get the VoteId of the first law
    voteURL = f"https://www.cdep.ro/pls/steno/evot2015.nominal?idv={vote_id}&idl=1"
    file_log("Fetching URL:", voteURL)
    prezentaKey = "prezenta"
    getlawForm = {
            #/pls/proiecte/upl_pck2015.proiect?\?idp=(\d+)
            VoteIdKey: r"/pls/proiecte/upl_pck2015.proiect?\?idp=(\d+)",
            prezentaKey: "Verificare prezenta"
        }
    results = fetch_and_search(voteURL, getlawForm)
    if results[prezentaKey] and results[VoteIdKey] is None:
        print("Special case: presence check only, no law associated.")
        file_log("Special case: presence check only, no law associated for URL:", voteURL)
        return {
            'titlu': "Verificare prezenta",
            'descriere': "Verificare prezenta in parlament pentru aceasta sesiune."
        }
    secondParams = {
            VoteIdKey: results[VoteIdKey] if results[VoteIdKey] else None
        }
    file_log("Got param law url:", secondParams)
    

    if not secondParams[VoteIdKey]:
        file_log("No VoteId found in the voting round page for URL:", voteURL)
        return None

    lawUrl=f"https://www.cdep.ro/pls/proiecte/upl_pck2015.proiect?idp={str(secondParams[VoteIdKey])}"
    file_log("Fetching law URL:", lawUrl)


    getlawForm = {
            # 1. EXPUNEREA DE MOTIVE
            # Example: /proiecte/2018/200/80/7/em361.pdf
            motivatieKey: re.compile(r'href=["\'](/proiecte/\d{4}/[\d/]+/em\d+\.pdf)["\']', re.I),
            # 2. FORMA ADOPTATĂ
            # Example: /proiecte/2018/200/80/7/se361.pdf
            # (Senat = se...pdf; Camera Deputaților often = cd...pdf)
            adoptedKey: re.compile(r'href=["\'](/proiecte/\d{4}/[\d/]+/(?:se|cd)\d+\.pdf)["\']', re.I),
            # 3. FORMA PENTRU PROMULGARE
            # Example: docs/2018/pr287_18.pdf
            promulgareKey: re.compile(r'href=["\'](docs/\d{4}/[^"\']*_pr[^"\']*\.pdf)["\']', re.I),
            # 4. ORDONANȚA DE URGENȚĂ
            # Example: /proiecte/2018/200/80/4/oug357.pdf
            ougKey: re.compile(r'href=["\'](/proiecte/\d{4}/[\d/]+/oug\d+\.pdf)["\']', re.I),
        }

    results = fetch_and_search(lawUrl, getlawForm)

    file_log("Got law douments:", results)


    text = getTextFromLawByForm(results)

    if text:
        file_log("Extracted text length:", len(text))
        if onlyDesc:
            prompt = f"Extrage o scurta descriere pentru legea urmatoare si returneaza-o in format JSON cu cheia descriere:\n\n{text[:5000]}"  # limit to first 4000 chars
        else:
            prompt = f"Extrage un titlu si o scurta descriere pentru legea urmatoare si returneaza-le in format JSON cu cheile titlu is descriere:\n\n{text[:5000]}"  # limit to first 4000 chars
        json_response = ask_chatgpt_json(prompt)
        file_log("ChatGPT JSON response:", json_response)
        return json_response
    else:
        file_log("No text extracted from the law document. For law URL:", lawUrl)

    return None


def tryPopulateAllLawsFromDBWithDefautlName():
    laws = getUnpopulatedLaws()
    for law in laws:
        time.sleep(2)  # to avoid hitting API rate limits
        try:
            onlyDesc = "Comprehensive legislation to reduce carbon emissions by" in str(law[DescriptionKey])
            if onlyDesc:
                file_log("Processing only description update for law with title:" + law[TitleKey], alsoPrint=True)
            if "Vot electronic" in str(law[TitleKey]) or "Comprehensive legislation to reduce carbon emissions by" in str(law[DescriptionKey]):
                file_log("Processing law with title:" + law[TitleKey], alsoPrint=True)
                desc = getNewNameAndDescriptionForLaw(law, onlyDesc=onlyDesc)
                if desc:
                    if onlyDesc:
                        file_log(f"Updating only description for law ID {law[IdKey]}.", alsoPrint=True)
                        update_law_description(law[IdKey], desc['descriere'])
                    else:
                        file_log(f"Updating TITLE and description for law ID {law[IdKey]}.", alsoPrint=True)
                        update_law(law[IdKey], desc)
                    file_log(f"Updated law ID {law.get(IdKey)} successfully.", alsoPrint=True)
                else:
                    file_log(f"No description produced for law ID {law.get(IdKey)}; deleting record.", alsoPrint=True)
                    try:
                        delete_law(law.get(IdKey))
                    except Exception as e:
                        file_log(f"Failed to delete law ID {law.get(IdKey)}: {e}", alsoPrint=True)
            else:
                file_log("Skipping law with title:" + law[TitleKey], alsoPrint=True)
        except Exception as e:
            file_log(f"Error processing law ID {law[IdKey]}: {e}", alsoPrint=True)        

if __name__ == "__main__":
    # Run the processing loop once every hour. Main-level prints/errors remain
    # on the console as requested; internal logs are written to output.txt.
    INTERVAL_SECONDS = 60 * 60  # 1 hour
    print("Starting hourly processing loop. Running once every hour.")
    try:
        while True:
            tryPopulateAllLawsFromDBWithDefautlName()

            # finished one pass
            file_log(f"Cycle complete — sleeping {INTERVAL_SECONDS} seconds until next run.", alsoPrint=True)
            time.sleep(INTERVAL_SECONDS)
    except KeyboardInterrupt:
        file_log("Interrupted by user; exiting.", alsoPrint=True)


        









