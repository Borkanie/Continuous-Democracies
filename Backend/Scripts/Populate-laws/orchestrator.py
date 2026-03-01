import sys
import os
import time
import datetime

SCRIPT_DIR = os.path.dirname(__file__)
if SCRIPT_DIR not in sys.path:
	sys.path.insert(0, SCRIPT_DIR)

from log_writer import log

# Direct imports (files live in this directory)
from import_politicians import importPoliticians
from import_law import goFromLastforward
from activate_politicians import activate_all_groups
from dataBase_interaction import init_db, close_db

def run_cycle():
	year = datetime.datetime.now().year
	log(f"Orchestrator cycle start: {datetime.datetime.now().isoformat()}")

	# 1) Import politicians only for the current year (IDs 1..333)
	try:
		log(f"Importing politicians for year {year}")
		importPoliticians(1, 333, year)
	except Exception as e:
		log(f"Failed to import politicians: {e}")

	# 2) Import new laws from last forward, attempt up to 100
	try:
		log("Calling import_law.goFromLastforward(100)")
		goFromLastforward(100)
	except Exception as e:
		log(f"Failed to import laws: {e}")

	log(f"Orchestrator cycle end: {datetime.datetime.now().isoformat()}")

def main():
	log("Starting orchestrator (runs once every hour).")

	# initialize shared DB connection once
	try:
		init_db()
		log("Database connection initialized.")
	except Exception as e:
		log(f"Failed to initialize DB connection: {e}")

	# trigger activation immediately on first loop by setting last_activate in the past
	last_activate = datetime.datetime.now() - datetime.timedelta(hours=12)

	try:
		while True:
			run_cycle()

			# every 12 hours run the activate script
			now = datetime.datetime.now()
			if (now - last_activate) >= datetime.timedelta(hours=12):
				try:
					log("Running activate_politicians.activate_all_groups()")
					activate_all_groups()
					last_activate = now
				except Exception as e:
					log(f"Failed to activate politicians: {e}")

			time.sleep(3600*12)
	except KeyboardInterrupt:
		log("Orchestrator stopped by user.")
	finally:
		try:
			close_db()
			log("Database connection closed.")
		except Exception:
			pass

if __name__ == "__main__":
	main()

