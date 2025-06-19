using ParliamentMonitor.Contracts.Model.Votes;

namespace ParliamentMonitor.WebInterface.ViewModels
{
    public class VoteContainer
    {
        public Vote Vote { get; set; } = null;

        public int PositionX = 0;
        public int PositionY = 0;

        private static void SetupLocationForContainers(ISet<VoteContainer> containerSet,int availableWidth, int availableHeigth)
        {
            double totalAreea = availableHeigth * availableHeigth;
            //double factorForoneSquareWithAvialableSpace = 2 * 3 * 4 * Math.Atan((double)(availableWidth) / (4 * availableHeigth)) / 10 ;
            // TODO check for inverse sizes here
            //double areeaOfOneSquare = factorForoneSquareWithAvialableSpace * totalAreea / containerSet.Count;
            double areeaOfOneSquare = 0.7 * totalAreea / containerSet.Count;

            int heigthOfSquare = (int)Math.Sqrt(areeaOfOneSquare * 4 / 3);
            int widthOfSquare = heigthOfSquare * 3 / 4;
            int currentRadius = availableHeigth * 2;

            int offsetX = 5;
            int offsetY = 5;
            foreach(var set in  containerSet)
            {
                set.PositionX = offsetX; 
                set.PositionY = offsetY;
                offsetX += widthOfSquare * 12 / 10;
                if(offsetX + widthOfSquare > availableWidth - 5)
                {
                    offsetX = 5;
                }
                offsetY += heigthOfSquare * 12 / 10;
                if(offsetY + heigthOfSquare > availableHeigth - 5)
                {
                    offsetY = 5;
                }
            }
        }
        
        public static ISet<VoteContainer> CreateContainers(ISet<Vote> votes, int width, int heigth)
        {
            var result = new HashSet<VoteContainer>();
            foreach(var vote in votes)
            {
                result.Add(new VoteContainer() { Vote = vote });
            }
            SetupLocationForContainers(result, width, heigth);
            return result;
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
