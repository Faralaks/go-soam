print("\n\n")

qests = open("qestions.txt").read().split("\n")
ans = open("answers.txt").read().split("\n")
results_file = open("result.txt", "w")

big_template = """
    <div class="card mb-5" id="card{}">
        <div class="card-header">Вопрос {}</div>
        <div class="card-body">
            <h5 class="card-title">{}</h5>

{}
            
        </div>
    </div>

"""


little_template = """            <div class="form-check">
                <input id="{}" class="form-check-input" type="radio" name="{}"  value="{}">
                <label for="{}" class="form-check-label">{}</label>
            </div>"""

for i, q in enumerate(qests, 1):
    variants = ""
    for j, a in enumerate(ans, 1):
        q_id = "q%da%d" % (i, j)
        variant = little_template.format(q_id, "q%d"%i, j, q_id, a).replace("\n\n", "") + "\n"
        variants += variant
    qestion = big_template.format(i, i, q, variants)
    results_file.write(qestion)


results_file.close()

