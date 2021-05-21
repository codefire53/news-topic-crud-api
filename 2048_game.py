import random
def put_random_init(cells):
    new_element = 0
    proba = random.uniform(0, 1)
    if (proba > 0.9):
        new_element = 4
    else:
        new_element = 2

    for i in range(2):
        r = random.randint(0, len(cells)-1)
        c = random.randint(0, len(cells[0])-1)
        while (cells[r][c] != 0):
            r = random.randint(0, len(cells)-1)
            c = random.randint(0, len(cells[0])-1)
        cells[r][c] = new_element
    return 

def check_conditions(cells, target):
    canContinue = False
    for i in range(len(cells)):
        for j in range(len(cells[i])):
            if cells[i][j]==target:
                return 'win'
            elif cells[i][j]==0:
                canContinue = True
            elif i < len(cells)-1 and cells[i][j]==cells[i+1][j]:
                canContinue = True
            elif j < len(cells[i])-1 and cells[i][j]==cells[i][j+1]:
                canContinue = True    
    if canContinue:
        return 'continue'
    return 'lose'

def put_random(cells):
    r = random.randint(0, len(cells)-1)
    c = random.randint(0, len(cells[0])-1)
    while (cells[r][c] != 0):
        r = random.randint(0, len(cells)-1)
        c = random.randint(0, len(cells[0])-1)
    cells[r][c] = 2
    return cells

def init_game(size):
    cells = []
    for i in range(size):
        cells.append([0]*size)
    put_random_init(cells)
    return cells

def slide_cells(cells):
    changed = False
    updated_cells = []
    for i in range(len(cells)):
        updated_cells.append([0]*len(cells[i]))
    for i in range(len(cells)):
        cur_col = 0
        for j in range(len(cells[i])):
            if(cells[i][j] != 0):
               updated_cells[i][cur_col] = cells[i][j]
               if (cur_col != j):
                   changed = True
               cur_col += 1
    return updated_cells, changed

def merge(cells):
    changed = False
    updated_cells = cells
    score = 0
    for i in range(len(updated_cells)):
        cur_col = 0
        for j in range(1, len(updated_cells[i])):
            if (updated_cells[i][j] == updated_cells[i][j-1] and updated_cells[i][j] != 0):
                updated_cells[i][j-1] = updated_cells[i][j-1] + updated_cells[i][j]
                score = updated_cells[i][j-1] + updated_cells[i][j]
                updated_cells[i][j] = 0
                changed = True
    return updated_cells, changed, score


def reverse(cells):
    for i in range(len(cells)):
        cells[i] = cells[i][::-1]
    return cells

def transpose(cells):
    updated_cells = []
    for i in range(len(cells)):
        updated_cells.append([])
        for j in range(len(cells[i])):
            updated_cells[i].append(cells[j][i])
    return updated_cells

def move_left(cells):
    increment = 0
    updated_cells, changed_pos = slide_cells(cells)
    updated_cells, changed_value, increment = merge(updated_cells)

    updated_cells, changed_pos2 = slide_cells(updated_cells)

    return updated_cells, changed_pos or changed_pos2 or changed_value, increment

def move_right(cells):
    increment = 0
    updated_cells=reverse(cells)
    updated_cells, changed, increment = move_left(updated_cells)
    updated_cells=reverse(updated_cells)
    return updated_cells, changed, increment

def move_up(cells):
    increment = 0
    updated_cells = transpose(cells)
    updated_cells, changed, increment = move_left(updated_cells)
    updated_cells = transpose(updated_cells)
    return updated_cells, changed, increment

def move_down(cells):
    increment = 0
    updated_cells = transpose(cells)
    updated_cells, changed, increment = move_right(updated_cells)
    updated_cells = transpose(updated_cells)
    return updated_cells, changed, increment

def print_cells(cells):
    for i in range(len(cells)):
        print(cells[i])
        
def random_decision(target):
    new_element = 0
    proba = random.uniform(0, 1)
    if (proba > 0.9):
        new_element = 4
    else:
        new_element = 2
    if (new_element == target):
        print("You Win!")
    else:
        print("Game Over!")

def is_power_two(num):
    if (num == 1):
        return True
    while (num != 1):
        if (num%2==1):
            return False
        num /= 2
    return True
if __name__ == '__main__':
    n = int(input("Insert cell size (nxn): "))
    if (n <= 0):
        raise Exception("N should be positive!")
    
    t = int(input("Insert target value: "))
    if (not is_power_two(t) or not t > 0):
        raise Exception("T should be power of two & positive!")
    if (n==1):
        random_decision(t)
        exit()
    print("Press w to slide up")
    print("Press s to slide down")
    print("Press a to slide left")
    print("Press d to side right")
    score = 0 
    cells = init_game(n)

while(True):
    print(f"Score: {score}")
    print(print_cells(cells))
    cmd = input("Insert direction: ")
    if (cmd == 'w'):
        cells, changed, increment = move_up(cells)
    elif (cmd == 's'):
        cells, changed, increment = move_down(cells)
    elif (cmd == 'a'):
        cells, changed, increment = move_left(cells)
    elif (cmd == 'd'):
        cells, changed, increment = move_right(cells)
    else:
        print("Input error! Please try again")
        continue
    if changed:
        score += increment
        state = check_conditions(cells,t)
        if ( state == 'win'):
            print("You Win!")
            break
        elif (state == 'lose'):
            print("Game Over")
            break
        else:
            put_random(cells)
            continue

