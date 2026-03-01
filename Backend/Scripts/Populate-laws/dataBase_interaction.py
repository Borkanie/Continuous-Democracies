import uuid
import psycopg2
from log_writer import log
from psycopg2 import sql
import os
from vote import Vote
from voting_rounds  import VotingRounds

promulgareKey = "promulgare"
adoptedKey = "adopted"
motivatieKey = "motivatie"
ougKey = "oug"

VoteIdKey = "VoteId"
DescriptionKey = "Description"
TitleKey = "Title"
IdKey = "Id"

# Database connection parameters
db_params = {
    'dbname': os.getenv('DBName'),
    'user': os.getenv('DBUser'),
    'password': os.getenv('DBPassword'),
    'host': '192.168.1.108',  # Change to your database host
    'port': '5432'        # Default PostgreSQL port
}

# Shared connection support: initialize once and reuse across calls.
_shared_connection = None
_shared_db_params = dict(db_params)

def init_db(params: dict | None = None):
    global _shared_connection, _shared_db_params
    if params:
        _shared_db_params = dict(params)
    if _shared_connection is None:
        _shared_connection = psycopg2.connect(**_shared_db_params)
    return _shared_connection

def _get_connection():
    global _shared_connection, _shared_db_params
    if _shared_connection is None:
        _shared_connection = psycopg2.connect(**_shared_db_params)
    try:
        if getattr(_shared_connection, 'closed', False):
            _shared_connection = psycopg2.connect(**_shared_db_params)
    except Exception:
        _shared_connection = psycopg2.connect(**_shared_db_params)
    return _shared_connection

def close_db():
    global _shared_connection
    try:
        if _shared_connection is not None and not getattr(_shared_connection, 'closed', False):
            _shared_connection.close()
    finally:
        _shared_connection = None

def getPoliticans() -> list[dict]:
    try:
        # Reuse shared connection
        connection = _get_connection()
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
        log(f"Retrieved {len(politicians)} politicians from the database.")
        return politicians
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        # Close the cursor only; shared connection is reused
        if cursor:
            cursor.close()

def getParties() -> list[dict]:
    try:
        connection = _get_connection()
        cursor = connection.cursor()

        # Define your SELECT query
        select_query = f"SELECT \"Id\", \"Acronym\", \"Name\"  FROM public.\"Parties\";"

        # Execute the query
        cursor.execute(select_query)

        parties = []
        # Iterate over the results and convert to a list of maps
        for row in cursor.fetchall():
            party = {"Id": row[0], "Acronym": row[1], "Name": row[2]}
            parties.append(party)
        log(f"Retrieved {len(parties)} parties from the database.")
        return parties
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        if cursor:
            cursor.close()

def get_max_vote_id() -> int | str | None:
    """
    Return the highest VoteId from public."VotingRounds".
    Tries numeric MAX(CAST(... AS BIGINT)) first, falls back to lexical ORDER BY DESC.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()

        # Try numeric max (works if VoteId is numeric or numeric-like)
        try:
            query = sql.SQL('SELECT MAX(CAST({voteid} AS BIGINT)) FROM public."VotingRounds";').format(
                voteid=sql.Identifier(VoteIdKey)
            )
            cursor.execute(query)
            res = cursor.fetchone()[0]
            if res is not None:
                log(f"Max numeric VoteId found: {res}")
                return int(res)
        except Exception:
            # Fallback to lexical max if casting fails
            query = sql.SQL('SELECT {voteid} FROM public."VotingRounds" ORDER BY {voteid} DESC LIMIT 1;').format(
                voteid=sql.Identifier(VoteIdKey)
            )
            cursor.execute(query)
            row = cursor.fetchone()
            if row:
                log(f"Max VoteId (fallback) found: {row[0]}")
                return row[0]
        return None

    except psycopg2.Error as e:
        print(f"Database error while fetching max VoteId: {e}")
        return None

    finally:
        if 'cursor' in locals() and cursor:
            cursor.close()
            
def get_last_vote_id() -> str | None:
    """
    Return the last (most-recent) VoteId from public."VotingRounds" using lexical ordering.

    This is a simple helper that performs:
      SELECT "VoteId" FROM public."VotingRounds" ORDER BY "VoteId" DESC LIMIT 1;

    Returns the VoteId as a string, or None if no rows exist or on error.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        query = sql.SQL('SELECT {voteid} FROM public."VotingRounds" ORDER BY {voteid} DESC LIMIT 1;').format(
            voteid=sql.Identifier(VoteIdKey)
        )
        cursor.execute(query)
        row = cursor.fetchone()
        if row:
            log(f"get_last_vote_id found: {row[0]}")
            return row[0]
        return None
    except psycopg2.Error as e:
        log(f"Database error while fetching last VoteId: {e}")
        return None
    finally:
        if 'cursor' in locals() and cursor:
            cursor.close()
            
def getUnpopulatedLaws() -> list[dict]:
    try:
        connection = _get_connection()
        cursor = connection.cursor()

        # Define your SELECT query (only rows whose Title contains "Vot electronic")
        select_query = f"SELECT \"{IdKey}\", \"{VoteIdKey}\", \"{TitleKey}\", \"{DescriptionKey}\" FROM public.\"VotingRounds\" WHERE \"{TitleKey}\" ILIKE '%Vot electronic%';"

        # Execute the query
        cursor.execute(select_query)

        laws = []
        # Iterate over the results and convert to a list of maps
        for row in cursor.fetchall():
            law = {IdKey: row[0], VoteIdKey: row[1], TitleKey: row[2], DescriptionKey: row[3]}
            laws.append(law)
        log(f"Retrieved {len(laws)} laws from the database.")
        return laws
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        if cursor:
            cursor.close()

def getLaws() -> list[dict]:
    try:
        connection = _get_connection()
        cursor = connection.cursor()

        # Define your SELECT query
        select_query = f"SELECT \"{IdKey}\", \"{VoteIdKey}\", \"{TitleKey}\", \"{DescriptionKey}\" FROM public.\"VotingRounds\";"

        # Execute the query
        cursor.execute(select_query)

        laws = []
        # Iterate over the results and convert to a list of maps
        for row in cursor.fetchall():
            law = {IdKey: row[0], VoteIdKey: row[1], TitleKey: row[2], DescriptionKey: row[3]}
            laws.append(law)
        log(f"Retrieved {len(laws)} laws from the database.")
        return laws
    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error: {e}")
        return []
    finally:
        if cursor:
            cursor.close()

def insertNewParty(acronym: str, partyName = "Unknown Party") -> str:
    """
    Insert a new party into the Parties table.

    Args:
        acronym: Acronym of the party to insert.

    Returns:
        The ID of the newly inserted party.
    """
    try:
        connection = _get_connection()
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
            partyName
        ))
        connection.commit()
        log(f"Inserted new party with ID {new_party_id} and acronym {acronym}.")

    except psycopg2.Error as e:
        # error prints should remain on console
        log(f"Database error while inserting Party: {e}")
        if connection:
            connection.rollback()
        new_party_id = None

    finally:
        if cursor:
            cursor.close()
    return new_party_id

def insertVotingRound(voting_round: VotingRounds):
    """
    Insert a new voting round into the VotingRounds table.

    Args:
        voting_round: VotingRounds object containing the data to insert.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        print(f"Inserting new voting round: {voting_round.__dict__}")
        log(f"Inserting new voting round: {voting_round.__dict__}")
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
        log(f"Inserted voting round ID {voting_round.Id} successfully.")

    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error while inserting voting round: {e}")
        if connection:
            connection.rollback()

    finally:
        if cursor:
            cursor.close()

def insertPolitician(name: str, partyId: str, imageUrl = "", gender = 2, active = True) -> str:
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
        connection = _get_connection()
        cursor = connection.cursor()
        log(f"Inserting new politician: Name={name}, PartyId={partyId}")
        # Prepare SQL insert statement safely using psycopg2.sql
        insert_query = sql.SQL("""
            INSERT INTO public."Politicians"
            ("Id", "Gender", "ImageUrl", "PartyId", "Active", "WorkLocation", "Name")
            VALUES (%s, %s, %s, %s, %s, %s, %s);
        """)
        id = str(uuid.uuid4())
        cursor.execute(insert_query, (
            id,
            gender,
            imageUrl,
            partyId,
            True,
            "0",
            name
        ))
        connection.commit()
        log(f"Inserted politician {name} with ID {id} successfully.")

    except psycopg2.Error as e:
        # error prints should remain on console
        log(f"Database error while inserting Politician: {e}")
        if connection:
            connection.rollback()
        id = None

    finally:
        if cursor:
            cursor.close()
    return id


def set_politician_active_by_name(name_like: str) -> int:
    """
    Set Active = true for politician rows whose Name matches the provided pattern (case-insensitive ILIKE).

    `name_like` should include SQL wildcards if desired (e.g. '%Popa Ştefan%').
    Returns the number of rows updated.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        log(f"Setting Active=true for politicians matching: {name_like}")
        update_query = sql.SQL('''
            UPDATE public."Politicians"
            SET "Active" = TRUE
            WHERE "Name" ILIKE %s
            RETURNING "Id";
        ''')
        cursor.execute(update_query, (name_like,))
        rows = cursor.fetchall()
        # commit to persist changes
        connection.commit()
        count = len(rows)
        log(f"Updated Active for {count} politician(s) matching '{name_like}'")
        return count
    except psycopg2.Error as e:
        log(f"Database error while updating politician Active flag: {e}")
        if connection:
            connection.rollback()
        return 0
    finally:
        if cursor:
            cursor.close()


def replace_active_politicians_by_names(names: list[str]) -> dict:
    """
    Replace the active politicians set in a single transaction.

    Steps (single transaction):
      1. Set ALL politicians Active = FALSE
      2. For each name in `names`, set Active = TRUE where Name ILIKE '%name%'

    Returns a dict with counts: {'deactivated': int, 'activated': int}
    """
    if names is None:
        names = []
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        log("Starting replace_active_politicians_by_names transaction")

        # Start transaction (psycopg2 starts one by default)
        cursor.execute('UPDATE public."Politicians" SET "Active" = FALSE;')
        deactivated = cursor.rowcount if cursor.rowcount is not None else 0

        activated_total = 0
        for raw_name in names:
            pattern = f"%{raw_name}%"
            cursor.execute('UPDATE public."Politicians" SET "Active" = TRUE WHERE "Name" ILIKE %s RETURNING "Id";', (pattern,))
            rows = cursor.fetchall()
            activated = len(rows)
            activated_total += activated
            log(f"Activated {activated} rows for pattern: {pattern}")

        connection.commit()
        log(f"replace_active_politicians_by_names committed: deactivated={deactivated}, activated={activated_total}")
        return {'deactivated': deactivated, 'activated': activated_total}

    except psycopg2.Error as e:
        log(f"Database error in replace_active_politicians_by_names: {e}")
        try:
            connection.rollback()
        except Exception:
            pass
        return {'deactivated': 0, 'activated': 0}
    finally:
        if cursor:
            cursor.close()

def update_law_description(law_id: int, description: str):
    """
    Update a law in the VotingRounds table by its ID using data from a JSON object.

    Args:
        law_id: ID of the law to update.
        law_json: JSON object with 'titlu' and 'descriere' keys.
        db_params: Dictionary with psycopg2 connection parameters.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        log(f"Updating law ID {law_id} with data: {description}")
        # Prepare SQL update statement safely using psycopg2.sql
        update_query = sql.SQL("""
            UPDATE public."VotingRounds"
            SET "Description" = %s
            WHERE "Id" = %s;
        """)


        cursor.execute(update_query, (description, law_id))
        connection.commit()
        log(f"Updated law ID {law_id} successfully.")

    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error while updating law ID {law_id}: {e}")
        if connection:
            connection.rollback()

    finally:
        if cursor:
            cursor.close()

def update_law(law_id: int, law_json: dict):
    """
    Update a law in the VotingRounds table by its ID using data from a JSON object.

    Args:
        law_id: ID of the law to update.
        law_json: JSON object with 'titlu' and 'descriere' keys.
        db_params: Dictionary with psycopg2 connection parameters.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        log(f"Updating law ID {law_id} with data: {law_json}")
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
        log(f"Updated law ID {law_id} successfully.")

    except psycopg2.Error as e:
        # error prints should remain on console
        print(f"Database error while updating law ID {law_id}: {e}")
        if connection:
            connection.rollback()

    finally:
        if cursor:
            cursor.close()


def delete_law(law_id: int):
    """
    Delete a law row from public."VotingRounds" by Id.

    Args:
        law_id: ID of the law to delete.
    """
    try:
        connection = _get_connection()
        cursor = connection.cursor()
        log(f"Deleting law ID {law_id} from database")
        delete_query = sql.SQL('''
            DELETE FROM public."VotingRounds"
            WHERE "Id" = %s;
        ''')
        cursor.execute(delete_query, (law_id,))
        connection.commit()
        log(f"Deleted law ID {law_id} successfully.")
    except psycopg2.Error as e:
        print(f"Database error while deleting law ID {law_id}: {e}")
        if connection:
            connection.rollback()
    finally:
        if cursor:
            cursor.close()


def insert_votes(votes: list[Vote]) -> int:
    """
    Insert multiple votes into public."Votes" table.

    Each vote dict must include the keys: "Id", "RoundId", "PoliticianId", "Position", "Name".

    Uses a single transaction and `ON CONFLICT ("Id") DO NOTHING` to skip duplicates if the
    table defines a unique constraint or primary key on "Id".

    Returns the number of rows successfully inserted (may be fewer than len(votes) if duplicates).
    """
    log(f"insert_votes: attempting to insert {len(votes)} votes into database.")
    if not votes:
        return 0

    inserted = 0
    try:
        connection = _get_connection()
        cursor = connection.cursor()

        # Prepare parameterized insert; use %s placeholders for psycopg2
        insert_sql = sql.SQL(
            'INSERT INTO public."Votes" ("Id", "RoundId", "PoliticianId", "Position", "Name") VALUES %s ON CONFLICT ("Id") DO NOTHING'
        )
        log("Prepared insert SQL for votes.")
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
            log("Using psycopg2.extras.execute_values for bulk insert.")
            execute_values(cursor, insert_sql.as_string(connection), values, template=None, page_size=100)
            # Get number of rows affected
            inserted += cursor.rowcount
            log(f"Inserted {inserted} votes.")
        else:
            log("psycopg2.extras.execute_values not available, falling back to single inserts.")
            # Fallback: insert rows one by one
            for row in values:
                try:
                    cursor.execute('INSERT INTO public."Votes" ("Id", "RoundId", "PoliticianId", "Position", "Name") VALUES (%s, %s, %s, %s, %s) ON CONFLICT ("Id") DO NOTHING', row)
                    # Get number of rows affected
                    inserted += cursor.rowcount
                    log(f"Inserted vote {row[0]} successfully.")
                except Exception as e:
                    log(f"Failed to insert vote {row[0]}: {e}")
        connection.commit()
        log("Committed vote inserts to database.")

    except psycopg2.Error as e:
        log(f"Database error while inserting votes: {e}")
        if connection:
            connection.rollback()
        return 0
    finally:
        if cursor:
            cursor.close()
    return inserted