class Vote:
    def __init__(self, votid, nume : str, prenume  :str, grup : str, camera, vot, roundId):
        self.id = None
        self.roundId = roundId
        self.politicianId = None
        self.partyId = None
        self.name = f"{prenume} {nume}"
        self.position = 0
        self.votid = votid
        self.grup = grup.strip()
        self.camera = camera
        self.vot = vot    
