function getInformation()
    name = "testItem"
    functionName = "rangeUp"
    numReturns = 1
    imageBounds = "0;142;22;157;"
    return imageBounds, numReturns, functionName, name
end

function rangeUp()
    rangeIncrease = 5

    SetGunRange(GunRange() + rangeIncrease)

    return true
end
