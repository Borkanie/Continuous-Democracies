import uuid
import psycopg2
from log_writer import file_log
from psycopg2 import sql
import os
from Vote import Vote
from VotingRounds  import VotingRounds

# Database connection parameters
db_params = {
    'dbname': os.getenv('DBName'),
    'user': os.getenv('DBUser'),
    'password': os.getenv('DBPassword'),
    'host': '192.168.1.108',  # Change to your database host
    'port': '5432'        # Default PostgreSQL port
}

def getPoliticans() -> list[dict]:
    try:
        # Connect to the PostgreSQL database
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()

        # Define your SELECT query
        select_query = f"SELECT \"Id\", \"PartyId\", \"Name\" FROM public.\"Politicians\";"

        # Execute the query
        cursor.execute(select_query)

        politicians = []
        # Iterate over the results and convert to a list of maps
        for row in cursor.fetchall():
            politician = {"Id": row[0], "PartyId": row[1], "Name": row[2]}
            politicians.append(politician)
        file_log(f"Retrieved {len(politicians)} politicians from the database.")
        return politicians
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        # Close the cursor and connection
        if cursor:
            cursor.close()
        if connection:
            connection.close()

def getParties() -> list[dict]:
    try:
        # Connect to the PostgreSQL database
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()

        # Define your SELECT query
        select_query = f"SELECT \"Id\", \"Acronym\"  FROM public.\"Parties\";"

        # Execute the query
        cursor.execute(select_query)

        parties = []
        # Iterate over the results and convert to a list of maps
        for row in cursor.fetchall():
            party = {"Id": row[0], "Acronym": row[1]}
            parties.append(party)
        file_log(f"Retrieved {len(parties)} parties from the database.")
        return parties
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        # Close the cursor and connection
        if cursor:
            cursor.close()
        if connection:
            connection.close()



def getLaws(IdKey : str, VoteIdKey : str, TitleKey : str) -> list[dict]:
    try:
        # Connect to the PostgreSQL database
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()

        # Define your SELECT query
        select_query = f"SELECT \"{IdKey}\", \"{VoteIdKey}\", \"{TitleKey}\"  FROM public.\"VotingRounds\";"

        # Execute the query
        cursor.execute(select_query)

        laws = []
        # Iterate over the results and convert to a list of maps
        for row in cursor.fetchall():
            law = {IdKey: row[0], VoteIdKey: row[1], TitleKey: row[2]}
            laws.append(law)
        file_log(f"Retrieved {len(laws)} laws from the database.")
        return laws
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        # Close the cursor and connection
        if cursor:
            cursor.close()
        if connection:
            connection.close()

def insertNewParty(acronym: str) -> str:
    """
    Insert a new party into the Parties table.

    Args:
        acronym: Acronym of the party to insert.

    Returns:
        The ID of the newly inserted party.
    """
    try:
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()
        insert_query = sql.SQL("""
            INSERT INTO public."Parties"("Id", "Acronym", "LogoUrl", "Color", "Active", "Name")
            VALUES (%s, %s, %s, %s, %s, %s);
        """)
        new_party_id = str(uuid.uuid4())
        cursor.execute(insert_query, (
            new_party_id,
            acronym,
            "",
            "#727272",
            True,
            "Unknown Party"
        ))
        connection.commit()
        file_log(f"Inserted new party with ID {new_party_id} and acronym {acronym}.")

    except psycopg2.Error as e:
        # error prints should remain on console
        file_log(f"Database error while inserting voting round: {e}", alsoPrint=True)
        if connection:
            connection.rollback()
        new_party_id = None

    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()
    return new_party_id

def insertVotingRound(voting_round: VotingRounds):
    """
    Insert a new voting round into the VotingRounds table.

    Args:
        voting_round: VotingRounds object containing the data to insert.
    """
    try:
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()
        print(f"Inserting new voting round: {voting_round.__dict__}")
        file_log(f"Inserting new voting round: {voting_round.__dict__}")
        # Prepare SQL insert statement safely using psycopg2.sql
        insert_query = sql.SQL("""
            INSERT INTO public."VotingRounds" ("Id", "VoteId", "Title", "VoteDate","Description","Name")
            VALUES (%s, %s, %s, %s, %s, %s);
        """)

        cursor.execute(insert_query, (
            voting_round.Id,
            voting_round.VoteId,
            voting_round.Title,
            voting_round.voteDate,
            "Default Description",
            voting_round.Title
        ))
        connection.commit()
        file_log(f"Inserted voting round ID {voting_round.Id} successfully.")

    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error while inserting voting round: {e}")
        if connection:
            connection.rollback()

    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()

def insertPolitician(name: str, partyId: str) -> str:
    """
    Insert a new politician record into the Politicians table.
    Args:
        name (str): The name of the politician to be inserted.
        partyId (str): The unique identifier of the party the politician belongs to.
    Returns:
        Id of the inserted Politician or None if there was an issue
    Raises:
        psycopg2.Error: If a database error occurs during insertion, the transaction is rolled back and an error message is printed.
    Side Effects:
        - Logs the insertion attempt to file and console.
        - Commits the transaction to the database on successful insertion.
        - Closes database cursor and connection resources.
    Notes:
        - A new UUID is generated as the politician's Id.
        - Gender is set to "2" by default.
        - ImageUrl is initialized as an empty string.
        - Active status is set to True.
        - WorkLocation is set to "0" by default.
    """
    
    try:
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()
        file_log(f"Inserting new politician: Name={name}, PartyId={partyId}", alsoPrint=True)
        # Prepare SQL insert statement safely using psycopg2.sql
        insert_query = sql.SQL("""
            INSERT INTO public."Politicians"
            ("Id", "Gender", "ImageUrl", "PartyId", "Active", "WorkLocation", "Name")
            VALUES (%s, %s, %s, %s, %s, %s, %s);
        """)
        id = str(uuid.uuid4())
        cursor.execute(insert_query, (
            id,
            "2",
            "",
            partyId,
            True,
            "0",
            name
        ))
        connection.commit()
        file_log(f"Inserted politician {name} with ID {id} successfully.", alsoPrint=True)

    except psycopg2.Error as e:
        # error prints should remain on console
        file_log(f"Database error while inserting voting round: {e}", alsoPrint=True)
        if connection:
            connection.rollback()
        id = None

    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()
    return id

def update_law(law_id: int, law_json: dict):
    """
    Update a law in the VotingRounds table by its ID using data from a JSON object.

    Args:
        law_id: ID of the law to update.
        law_json: JSON object with 'titlu' and 'descriere' keys.
        db_params: Dictionary with psycopg2 connection parameters.
    """
    try:
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()
        file_log(f"Updating law ID {law_id} with data: {law_json}")
        # Prepare SQL update statement safely using psycopg2.sql
        update_query = sql.SQL("""
            UPDATE public."VotingRounds"
            SET "Title" = %s,
                "Description" = %s
            WHERE "Id" = %s;
        """)

        title = law_json.get("titlu")
        description = law_json.get("descriere")

        cursor.execute(update_query, (title, description, law_id))
        connection.commit()
        file_log(f"Updated law ID {law_id} successfully.")

    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error while updating law ID {law_id}: {e}")
        if connection:
            connection.rollback()

    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()


def delete_law(law_id: int):
    """
    Delete a law row from public."VotingRounds" by Id.

    Args:
        law_id: ID of the law to delete.
    """
    try:
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()
        file_log(f"Deleting law ID {law_id} from database")
        delete_query = sql.SQL('''
            DELETE FROM public."VotingRounds"
            WHERE "Id" = %s;
        ''')
        cursor.execute(delete_query, (law_id,))
        connection.commit()
        file_log(f"Deleted law ID {law_id} successfully.")
    except psycopg2.Error as e:
        print(f"Database error while deleting law ID {law_id}: {e}")
        if connection:
            connection.rollback()
    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()


def insert_votes(votes: list[Vote]) -> int:
    """
    Insert multiple votes into public."Votes" table.

    Each vote dict must include the keys: "Id", "RoundId", "PoliticianId", "Position", "Name".

    Uses a single transaction and `ON CONFLICT ("Id") DO NOTHING` to skip duplicates if the
    table defines a unique constraint or primary key on "Id".

    Returns the number of rows successfully inserted (may be fewer than len(votes) if duplicates).
    """
    file_log(f"insert_votes: attempting to insert {len(votes)} votes into database.", alsoPrint=True)
    if not votes:
        return 0

    inserted = 0
    try:
        connection = psycopg2.connect(**db_params)
        cursor = connection.cursor()

        # Prepare parameterized insert; use %s placeholders for psycopg2
        insert_sql = sql.SQL(
            'INSERT INTO public."Votes" ("Id", "RoundId", "PoliticianId", "Position", "Name") VALUES %s ON CONFLICT ("Id") DO NOTHING'
        )
        file_log("Prepared insert SQL for votes.")
        # Build list of tuples for execute_values
        values = []
        for v in votes:
            values.append((v.id, v.roundId, v.politicianId, v.position, v.name))

        # Use psycopg2.extras.execute_values for efficient multi-row insert
        try:
            from psycopg2.extras import execute_values
        except Exception:
            execute_values = None
        
        if execute_values:
            file_log("Using psycopg2.extras.execute_values for bulk insert.", alsoPrint=True)
            execute_values(cursor, insert_sql.as_string(connection), values, template=None, page_size=100)
            # Get number of rows affected
            inserted += cursor.rowcount
            file_log(f"Inserted {inserted} votes.", alsoPrint=True)
        else:
            file_log("psycopg2.extras.execute_values not available, falling back to single inserts.", alsoPrint=True)
            # Fallback: insert rows one by one
            for row in values:
                try:
                    cursor.execute('INSERT INTO public."Votes" ("Id", "RoundId", "PoliticianId", "Position", "Name") VALUES (%s, %s, %s, %s, %s) ON CONFLICT ("Id") DO NOTHING', row)
                    # Get number of rows affected
                    inserted += cursor.rowcount
                    file_log(f"Inserted vote {row[0]} successfully.", alsoPrint=True)
                except Exception as e:
                    file_log(f"Failed to insert vote {row[0]}: {e}")
        connection.commit()
        file_log("Committed vote inserts to database.")

    except psycopg2.Error as e:
        file_log(f"Database error while inserting votes: {e}", alsoPrint=True)
        if connection:
            connection.rollback()
        return 0
    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()
    return inserted