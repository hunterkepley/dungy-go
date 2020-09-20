function getInformation()
    name = "healthPack"
    functionName = "heal"
    numReturns = 1
    imageBounds = "0;96;15;105;"
    return imageBounds, numReturns, functionName, name
end

function heal()
    healRate = 2
    
    SetPlayerHealth(PlayerHealth() + healRate)
    
    return true
end
