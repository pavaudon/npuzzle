\\\ MARCHE A SUIVRE sur l'interface pour les bugs ///


RESOLU 
_ Envoyer a la main la meme map
    => tourne en boucle

RESOLU
_ Envoyer la map :
    DonTKnow/
        3-11-moves.txt > IL FAUT QUE TU essaies de nouveau, en snail ?
        4.txt > La aussi je ne sais pas, voir le fichier, mais meme a resoudre a la main c'est compliqué, je pense qu'elle n'est pas resoluble
        npuzzle-3-1.txt >  > PAS DE SOUCIS
        npuzzle-4-1.txt > Voir le com du fichier
        solvable.txt > PAS DE SOUCIS
    ORDERED/
        5.1.txt >> elle met presque 2 min , c'est une 5x5 > 682ms
    SNAIL/
        3.1.txt >> Aucun soucis
    => tourne en boucle sans messages de resolution impossible (avec mon pc sur le point d'exploser)
    -> SOIT ca marche chez toi parce que tu as un meilleur pc et dans ce cas il faut voir si les macs de 42 sont aussi puissants
    -> SOIT il y a vraiment un probleme


/// POSSIBILITES D'AJOUTS?! \\\


TO PAULINE :
- Tu pourras verifier que closed = complexity time et Max open = complexity size?
- Les conditions de verification des map bonnes ou non, elles sont definies comment ?

||| BUG EN + |||
RESOLU
- Une map deja triée en snail = [ERROR] Map is not solvable
- attendre la fin de l'impression des MOVES avant d'exit le programme
  RESOLU
  _ 'start' pour une map deja triee et cocher 'annuler' (revenir a la map initiale pas triee)
    => 'UNABLE MOVE'
        -> map mal triee
        -> je pense que ce qu'on devrait faire c'est surement un reset de nos data par exemple de open/closed
