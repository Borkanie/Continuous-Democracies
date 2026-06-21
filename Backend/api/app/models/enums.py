import enum


class Gender(int, enum.Enum):
    Male = 0
    Female = 1
    Other = 2


class WorkLocation(int, enum.Enum):
    Parliament = 0
    Senate = 1
    Other = 2


class VotePosition(int, enum.Enum):
    Yes = 0
    No = 1
    Abstain = 2
    Absent = 3
