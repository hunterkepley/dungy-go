function getInformation()
    name = "testItem"
    functionName = "test"
    numReturns = 1
    imageBounds = "0;0;15;5;"
    return imageBounds, numReturns, functionName, name
end

function test()
    print("test")
    return true
end