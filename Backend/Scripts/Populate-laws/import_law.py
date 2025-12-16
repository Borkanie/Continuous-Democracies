import os
from log_writer import file_log
from utils import get, extract_datetime_from_dom, createtempDir, removeTempDir
from dataBaseInteraction import *
import xml.etree.ElementTree as ET
import uuid
from Vote import Vote
from VotingRounds  import VotingRounds

def extractNewVotingRound(lawId: int) -> VotingRounds:
    """
    Extract and store a new voting round for a given law.
    This function retrieves voting data from the Romanian Parliament website,
    parses the voting round information, and adds it to the database.
    Args:
        lawId (int): The unique identifier of the law for which to extract the voting round.
    Returns:
        VotingRounds: A VotingRounds object containing the extracted voting round data,
                      including a unique ID, the law ID, and the datetime of the vote.
    Raises:
        Potentially raises exceptions from get() or addNewVoteingRound() if HTTP request
        or database operations fail.
    Note:
        - Logs extraction activity to both console and file for audit purposes.
        - Makes an HTTP GET request to the Romanian Parliament website (cdep.ro).
        - The voting round is assigned a unique UUID upon creation.
    """

    print(f"Extracting new voting round for law ID: {lawId}")
    file_log(f"Extracting new voting round for law ID: {lawId}")
    # Implementation to extract and add a new voting round to the database
    getUrl = f"https://cdep.ro/pls/steno/evot2015.Nominal?idv={lawId}"
    response_text = get(url=getUrl, file=f"voting_round_{lawId}.html")
    round = VotingRounds(
        str(uuid.uuid4()),
        lawId,
        extract_datetime_from_dom(response_text)
        )
    insertVotingRound(round)
    return round

def addPoliticianMatching(politicians, vote : Vote) -> bool:
    matching_politician = [politician for politician in politicians if politician["Name"] == vote.name]

    if len(matching_politician) > 0:
        matching_politician = matching_politician[0]
        file_log(f"Matched politician: {matching_politician}")
        vote.id = str(uuid.uuid4())
        vote.politicianId = matching_politician["Id"]
    else:
        file_log(f"No matching politician found for vote: {vote.name}")
        return insertPolitician(vote.name, vote.partyId) is not None

    return True

def findPartyForVote(parties, vote : Vote) -> str:
    if vote.grup == "Minoritati":
        vote.grup = "MIN"
    matching_parties = [party for party in parties if party["Acronym"] == vote.grup]
    if len(matching_parties) == 0:
        file_log(f"No matching party found for acronym: {vote.grup}")
        id = insertNewParty(vote.grup)
        if id is None:
            file_log(f"Failed to insert new party for acronym: {vote.grup}", alsoPrint=True)
            return False
        else:
            file_log(f"Inserted new party with ID: {id} for acronym: {vote.grup}", alsoPrint=True)
        vote.partyId = id
    else:
        file_log(f"Matched party: {matching_parties}")
        vote.partyId = matching_parties[0]["Id"]
    return True
    
def setVotePosition(vote: Vote) -> None:
    """
    Sets the position attribute of a vote based on its vote value.
    Args:
        vote (Vote): The vote object to update.
    """
    if vote.vot == "DA":
        vote.position = 0
    elif vote.vot == "NU":
        vote.position = 1
    elif vote.vot == "AB":
        vote.position = 2
    else:
        vote.position = 3 

def getVoteListForALawById(lawId : int, existing_law_id : str) -> list[Vote]:
    """
        Retrieves and processes voting data for a specific law from the Romanian Chamber of Deputies.
        Args:
            lawId (int): The identifier of the law to retrieve votes for.
            existing_law_id (str): The internal identifier for the law in the local system.
        Returns:
            list[Vote]: A list of Vote objects, each representing a processed vote with matched politician and party information.
        Side Effects:
            - Logs processing steps and matches to a file.
            - Adds new politicians to the system if they are not found in the existing list.
        Notes:
            - Only votes with both a matched politician and party are included in the returned list.
            - The function assigns a position to each vote based on the vote value ("DA", "NU", "AB", or others).
    """

    url = f"https://www.cdep.ro/pls/steno/evot2015.xml?par1=2&par2={lawId}"
    response_text = get(url, file=f"law_{lawId}_response.txt")
    
    root = ET.fromstring(response_text)
    votes = [Vote(
        row.find('VOTID').text,
        row.find('NUME').text,
        row.find('PRENUME').text,
        row.find('GRUP').text,
        row.find('CAMERA').text,
        row.find('VOT').text,
        existing_law_id
    ) for row in root.findall('ROW')]
    print(f"Processed {len(votes)} votes for law ID {lawId}")
    parties = getParties()
    politicians= getPoliticans()
    for vote in votes:
        file_log(f"Processing vote for politician: {vote.name}, group: {vote.grup}")
        if not findPartyForVote(parties, vote):
            file_log(f"Skipping vote for politician: {vote.name} due to missing party.", alsoPrint=True)
            continue
        if not addPoliticianMatching(politicians, vote):
            file_log(f"Skipping vote for politician: {vote.name} due to missing politician.", alsoPrint=True)
            continue
        setVotePosition(vote)
    result = [vote for vote in votes if vote.politicianId is not None and vote.partyId is not None]
    file_log(f"Total valid votes after matching: {len(result)}", alsoPrint=True)
    return result


def checkifLawIsNotEMpty(lawId: int, file:str) -> bool:
    url = f"https://www.cdep.ro/pls/steno/evot2015.xml?par1=2&par2={lawId}"
    response_text = get(url, file=file, keep=True)

    root = ET.fromstring(response_text)
    votes = [row for row in root.findall('ROW')]
    return len(votes) > 0

def importLaws(startingIndex : int, endingIndex : int) -> list:
    existinglaws = getLaws("Id", "VoteId", "Title")
    for lawId in range(startingIndex, endingIndex + 1):
        try:
            file = f"law_{lawId}_response.txt"
            if checkifLawIsNotEMpty(lawId, file) == False:
                file_log(f"Law ID {lawId} has no votes. Skipping import.", alsoPrint=True)
                try:
                    os.remove(os.path.join(createtempDir(), file))
                except Exception:
                    file_log(f"Failed to remove temp file for law ID {lawId}: {file}")
                continue
            existing_law_id = None
            if any(law for law in existinglaws if law["VoteId"] == lawId):
                file_log(f"Law ID {lawId} already exists in the database. Skipping import.")
                existing_law_id = next(law for law in existinglaws if law["VoteId"] == lawId)["Id"]
            else:
                existing_law_id = extractNewVotingRound(lawId).Id

            votes = getVoteListForALawById(lawId, existing_law_id)
            file_log(f"Imported {len(votes)} votes for law ID {lawId}", alsoPrint=True)
            if insert_votes(votes) == 0:
                file_log(f"No votes were inserted for law ID {lawId}", alsoPrint=True)
            
        except Exception as e:
            file_log(f"Error processing law ID {lawId}: {e}")



if __name__ == "__main__":
    createtempDir()
    importLaws(250, 2000)
    removeTempDir()
    