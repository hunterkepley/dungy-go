function getInformation()
    name = "cacto"
    functionName = "thorny"
    numReturns = 1
    imageBounds = "0;114;15;129;"
    return imageBounds, numReturns, functionName, name
end

function thorny() -- makes the player run quickly for a few seconds in return for health
    
    damageAmount = 1
    speedIncrease = 1

    SetPlayerHealth(PlayerHealth()-damageAmount)

    SetPlayerWalkSpeed(PlayerWalkSpeed()+speedIncrease)
    SetPlayerRunSpeed(PlayerRunSpeed()+speedIncrease*2)

    return true
end
