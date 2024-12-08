# Sorted Map

A sorted map implementation. Probably just one of the millions out there.  It uses a map, an `RWMutex` and a sorted slice underneath. It also uses binary search for insertion and deletion of the keys to make it performant.
