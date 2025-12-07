import psycopg2
from log_writer import file_log
from psycopg2 import sql

# Database connection parameters
db_params = {
    'dbname': 'parlimentdb',
    'user': 'bobo',
    'password': 'password123',
    'host': '192.168.1.108',  # Change to your database host
    'port': '5432'        # Default PostgreSQL port
}


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