function getInformation()
    name = "brain"
    functionName = "smarter"
    numReturns = 1
    imageBounds = "0;129;21;143;"
    return imageBounds, numReturns, functionName, name
end

function smarter() -- makes the player have thorns when touching enemies
    accuracyChange = 40

    if Accuracy() < 100 then
        SetAccuracy(Accuracy() + accuracyChange)
    end
    
    return true
end
