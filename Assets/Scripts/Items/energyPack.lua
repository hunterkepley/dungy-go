function getInformation()
    name = "energyPack"
    functionName = "energize"
    numReturns = 1
    imageBounds = "0;105;15;114;"
    return imageBounds, numReturns, functionName, name
end

function energize()
    energizeRate = 4
    
    SetPlayerEnergy(PlayerEnergy() + energizeRate)
    
    return true
end
