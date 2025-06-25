using ParliamentMonitor.Contracts.Model.Votes;

namespace ParliamentMonitor.WebInterface.ViewModels
{
    public class VoteContainer
    {
        public Vote Vote { get; set; } = null;

        public double PositionX = 0;
        public double PositionY = 0;

        public static int VoteBoxWidth { get; private set; } = 0;
        public static int VoteBoxHeigth { get; private set; } = 0;

        private static void SetupLocationForContainers(IList<VoteContainer> containerSet,int availableWidth, int availableHeigth, int topX = 0, int topY = 0)
        {
            double totalAreea = availableHeigth * availableHeigth;
            double areeaOfOneSquare = 0.7 * totalAreea / containerSet.Count;

            VoteBoxHeigth = (int)Math.Sqrt(areeaOfOneSquare * 4 / 3);
            VoteBoxWidth = VoteBoxHeigth * 3 / 4;


            int centralPointforCirclesX = availableWidth / 2;
            int centralPointforCircleY = availableHeigth * 7 / 8;

            double currentAngle = 0;
            int radiusIncrease = 50;
            int radius = 150;
            double numberOfPlaces = 7;
            double currentPlace = 0;
            for(int i=0;i< containerSet.Count;i++)
            {
                var container = containerSet.ElementAt(i);
                container.PositionX = topX + centralPointforCirclesX - Math.Cos(currentAngle) * radius;
                container.PositionY = topY + centralPointforCircleY - Math.Sin(currentAngle) * radius;

                currentAngle = currentPlace / numberOfPlaces * Math.PI;
                currentPlace++;
                if(currentPlace > numberOfPlaces)
                {
                    radius += radiusIncrease;
                    numberOfPlaces += 7;
                    currentPlace = 0;
                    if( numberOfPlaces > containerSet.Count - i)
                    {
                        numberOfPlaces = containerSet.Count - i;
                    }
                }
            }
        }
        
        public static ISet<VoteContainer> CreateContainers(ISet<Vote> votes, int width, int heigth)
        {
            var result = new List<VoteContainer>();
            foreach(var vote in votes)
            {
                result.Add(new VoteContainer() { Vote = vote });
            }
            SetupLocationForContainers(result, width, heigth);
            return result.ToHashSet();
        }

        public override bool Equals(object? obj)
        {
            if (obj is not VoteContainer)
                return false;
            var ob = (VoteContainer)obj;
            if(ob.Vote!= Vote || ob.PositionY != PositionY || ob.PositionX!= PositionX) return false;
            return true;
        }
    }
}
